import path from "path";
import yauzl from "yauzl";
import { cwd } from "process";
import fs from "fs";
import { FileNotFoundError } from "./fileError.js";

export class ServerFileService {
    private static _instance: ServerFileService;
    private path: string;
    // Problem: this.path relies on serverHandler, however, serverHandler will likely call this calss on initialization

    private constructor() {
        this.path = cwd();
        // TODO: use alternative if mcdr exists
    }

    public static getInstance(): ServerFileService {
        if (!ServerFileService._instance)
            ServerFileService._instance = new ServerFileService();
        return ServerFileService._instance;
    }

    public async readFileFromZip(zip: string, target: string): Promise<string> {
        // Use this method for jar files
        if (!this.exist(zip)) throw FileNotFoundError;
        const zipPath = path.join(cwd(), zip);
        return new Promise((resolve, reject) => {
            yauzl.open(zipPath, { lazyEntries: true }, (err, zipFile) => {
                if (err) return reject(err);
                zipFile.on("entry", (entry: yauzl.Entry) => {
                    if (entry.fileName === target) {
                        zipFile.openReadStream(entry, (err, readStream) => {
                            if (err) return reject(err);
                            const chunks: Buffer[] = [];
                            readStream.on("data", (chunk) => {
                                chunks.push(chunk);
                            });
                            readStream.on("end", () => {
                                return resolve(
                                    Buffer.concat(chunks).toString()
                                );
                            });
                            readStream.on("error", reject);
                        });
                    } else {
                        zipFile.readEntry();
                    }
                });
                zipFile.readEntry();
            });
        });
    }

    public async readFrom(target: string): Promise<string> {
        if (!this.exist(target)) throw FileNotFoundError;
        const targetPath = path.join(cwd(), target);
        return new Promise((resolve, reject) => {
            fs.readFile(targetPath, "utf-8", (err, data) => {
                if (err) return reject(err);
                return resolve(data);
            });
        });
    }

    public exist(...targets: string[]): boolean {
        return targets.every((target) =>
            fs.existsSync(path.join(cwd(), target))
        );
    }
}
