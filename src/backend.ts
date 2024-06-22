import fs from "fs";
import path from "path";
import yauzl from "yauzl";
import { Server } from "./probe.js";

// Service logic
class Service {
    private generalOperations(): void {}

    private commandInit(): void {
        fs.mkdirSync(".mcpm");
        this.probeServer();
    }

    private probeServer(): Server {
        let server: Server = new Server();
        return server;
    }
}

// File related jobs
class FileSystem {
    private pwd: string = process.cwd();
    private mcpmPath = ".mcpm";

    public async readJsonFromJar(
        jarPath: string, // This should be relative to pwd
        target: string // This should be relative to the .jar file
    ): Promise<JSON> {
        if (/\.json$/i.test(target) === false) {
            throw "Param 'target' do not point to a JSON file.";
        }

        const extractPath = path.join(this.pwd, this.mcpmPath);
        const jarStream = fs.createReadStream(jarPath);

        yauzl.open(jarPath, (err, zipfile) => {
            if (err) throw err;
            zipfile.on("read", (stream) => {});
        });

        return JSON.parse("");
    }

    // public deleteFile() {}
    // public writeToFile(){}
}

// Network related jobs: fetch api, download, etc.
class Networking {}

export { Service, FileSystem, Networking };
