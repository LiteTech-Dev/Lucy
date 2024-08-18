import { cwd } from "process";
// import properties from "properties-reader";

import { Mcdr, McdrConfig, Minecraft, ModLoader } from "./classServerModule.js";
import { readMcdrConfig, searchForServerExecutable } from "./serverProbing.js";

export class ServerHandler {
    private static _instance: ServerHandler;
    private modLoader?: ModLoader;
    private mcdr?: Mcdr;
    private minecraft: Minecraft = new Minecraft("1.0");
    path: string = cwd();
    executable: string = "";

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
        try {
            this.executable = searchForServerExecutable(this.path);
        } catch (err) {
            // TODO: Return this to the user
            // If no executable is found, exit and notify the user
            // If multiple executables are found, let them choose a default
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
