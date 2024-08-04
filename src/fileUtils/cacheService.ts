import fs from "fs";
import { tmpdir } from "os";
import path from "path";

// Do not use this yet, it is unsure if the caching is necessary

export class CacheService {
    private static _instance: CacheService;
    private readonly cachePath = path.join(tmpdir(), "lucy");

    private constructor() {}

    public static getInstance(): CacheService {
        if (!CacheService._instance)
            CacheService._instance = new CacheService();
        return CacheService._instance;
    }

    public createCacheDirectory(): void {
        if (!fs.existsSync(this.cachePath)) {
            fs.mkdirSync(this.cachePath);
        }
    }

    public async readFrom(target: string): Promise<string> {
        const targetPath = path.join(this.cachePath, target);
        return new Promise((resolve, reject) => {
            fs.readFile(targetPath, "utf-8", (err, data) => {
                if (err) reject(err);
                resolve(data);
            });
        });
    }

    public async writeTo(target: string, data: string): Promise<void> {
        const targetPath = path.join(this.cachePath, target);
        return new Promise((resolve, reject) => {
            fs.writeFile(targetPath, data, (err) => {
                if (err) reject(err);
                resolve();
            });
        });
    }

    public async clearCache(): Promise<void> {
        // A delete method was not implemented because no specific file should be deleted from the cache
        return new Promise((resolve, reject) => {
            fs.rmdir(this.cachePath, { recursive: true }, (err) => {
                if (err) reject(err);
                resolve();
            });
        });
    }
}
