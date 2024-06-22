import fs from "fs";
import path from "path";
import yauzl from "yauzl";
import { Server } from "./server-modules.js";

const mcpmPath = ".mcpm";

// Service logic
class Service {
    fileSystem: FileSystem = new FileSystem();
    networking: Networking = new Networking();
    server: Server = this.probeServer();

    private probeServer(): Server {
        let server: Server = new Server();
        return server;
    }
}

// File related jobs
class FileSystem {
    private pwd: string = process.cwd();

    public async readJsonFromJar(
        jarPath: string, // This should be relative to pwd
        target: string // This should be relative to the .jar file
    ): Promise<JSON> {
        if (/\.json$/i.test(target) === false) {
            throw "Param 'target' do not point to a JSON file.";
        }

        const extractPath = path.join(this.pwd, mcpmPath);
        const jarStream = fs.createReadStream(jarPath);

        yauzl.open(jarPath, (err, zipfile) => {
            if (err) throw err;
            zipfile.on("entry", (stream) => {});
        });

        return JSON.parse("");
    }

    // public deleteFile() {}
    // public writeToFile(){}

    constructor() {
        if (!fs.existsSync(mcpmPath)) {
            fs.mkdirSync(mcpmPath);
        }
    }
}

// Network related jobs: fetch api, download, etc.
class Networking {}

export { Service };
