import { cwd } from "process";
// import properties from "properties-reader";

import {
    Mcdr,
    McdrConfig,
    Minecraft,
    ModLoader,
    ServerExecutable,
} from "./classServerModule.js";
import { readMcdrConfig, searchForServerExecutable } from "./serverProbing.js";
import {
    MultipleExecutablesFoundError,
    NoExecutableFoundError,
} from "./serverError.js";

export class ServerHandler {
    private static _instance: ServerHandler;
    private modLoader?: ModLoader;
    private mcdr?: Mcdr;
    private minecraft: Minecraft = new Minecraft("1.0");
    path: string = cwd();
    executable: ServerExecutable;

    private constructor() {
        // Initialize MCDR
        // TODO: Catch error from readMcdrConfig()
        const mcdrConfig: McdrConfig = readMcdrConfig();

        if (mcdrConfig) {
            const mcdr = new Mcdr(mcdrConfig);
            this.mcdr = mcdr;
            this.path = mcdr.mcdrServerPath;
        }

        // Search for server jar
        const validExecutables = searchForServerExecutable(this.path);
        switch (validExecutables.length) {
            case 0:
                throw NoExecutableFoundError;

            case 1:
                this.executable = validExecutables[0];
                break;

            default:
                throw MultipleExecutablesFoundError;
        }

        // Analyze the executable to initialize minecraft and modLoader
    }

    // Global access point
    public static getInstance(): ServerHandler {
        if (!ServerHandler._instance) {
            ServerHandler._instance = new ServerHandler();
        }
        return ServerHandler._instance;
    }

    public getModList() {}

    public addMod() {}

    public removeMod() {}

    public disableMod() {}

    public enableMod() {}
}
