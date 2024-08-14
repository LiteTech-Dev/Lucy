import { cwd } from "process";
import path from "path";
import fs from "fs";
import jsyaml from "js-yaml";
import properties from "properties-reader";

import { ServerFileService } from "../fileUtils/serverFileService.js";
import { ServerModule, Mcdr, McdrConfig } from "./classServerModule.js";

class ServerExecutableError extends Error {
    constructor(message: string) {
        super(message);
        this.name = "ServerExecutableError";
    }
}

export class Server {
    private static _instance: Server;
    path: string;
    private modules: ServerModule[] = [];

    private constructor() {
        // Initialize MCDR
        const mcdrConfig: McdrConfig = this.readMcdrConfig();

        if (mcdrConfig) {
            const mcdr = new Mcdr(mcdrConfig);
            this.modules.push(mcdr);
            this.path = mcdr.mcdrServerPath;
        } else {
            this.path = cwd();
        }

        // Search for server jar
    }

    // Global access point
    public static getInstance(): Server {
        if (!Server._instance) {
            Server._instance = new Server();
        }
        return Server._instance;
    }

    public moduleExists(moduleName: string): boolean {
        return this.modules.some((module) => module.name === moduleName);
    }

    private getModule(moduleName: string): ServerModule | null {
        return (
            this.modules.find((module) => module.name === moduleName) || null
        );
    }





    }

