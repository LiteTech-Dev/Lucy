import path from "path";
import fs from "fs";
import { cwd } from "process";

export class AppDataService {
    private static _instance: AppDataService;
    private readonly appDataPath = path.join(cwd(), ".lucy");

    private constructor() {}

    public static getInstance(): AppDataService {
        if (!AppDataService._instance)
            AppDataService._instance = new AppDataService();
        return AppDataService._instance;
    }

    public createAppDataPath(): void {
        if (!fs.existsSync(this.appDataPath)) {
            fs.mkdirSync(this.appDataPath);
        }
    }

    public appDataPathExists(): boolean {
        return fs.existsSync(this.appDataPath);
    }

    public exist(...targets: string[]): boolean {
        return targets.every((target) =>
            fs.existsSync(path.join(cwd(), target))
        );
    }

    public async readFrom(target: string): Promise<string> {
        if (!this.exist(target)) return "";
        const targetPath = path.join(this.appDataPath, target);
        return new Promise((resolve, reject) => {
            fs.readFile(targetPath, "utf-8", (err, data) => {
                if (err) return reject(err);
                resolve(data);
            });
        });
    }

    public async writeTo(target: string, data: string): Promise<void> {
        // This method overwrites the file if it already exists
        const targetPath = path.join(this.appDataPath, target);
        return new Promise((resolve, reject) => {
            fs.writeFile(targetPath, data, (err) => {
                if (err) return reject(err);
                resolve();
            });
        });
    }

    public async delete(target: string): Promise<void> {
        if (!this.exist(target)) return;
        const targetPath = path.join(this.appDataPath, target);
        return new Promise((resolve, reject) => {
            fs.unlink(targetPath, (err) => {
                if (err) return reject(err);
                resolve();
            });
        });
    }
}
