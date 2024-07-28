import { cwd } from "process";
import path from "path";
import jsyaml from "js-yaml";
import fs from "fs";

export class Server {
    private static _instance: Server;
    path: string = cwd();

    private constructor() {}

    // Global access point
    public static getInstance(): Server {
        if (!Server._instance) {
            Server._instance = new Server();
        }
        return Server._instance;
    }

    // addModule(module: ServerModule) {
    //     this[module.name] = module;
    // }
}

class ServerModule {
    name!: string;
}

interface McdrConfig {
    working_directory: string;
    plugin_directories: string[];
}

// Not tested yet
function getMcdr(this: Server): ServerModule {
    const mcdr = new ServerModule() as ServerModule & {
        mcdrServerPath: string;
        mcdrPluginsPath: string[];
    };
    mcdr.name = "MCDR";
    const mcdrConfigYml: string = fs
        .readFileSync(path.join(cwd(), "config.yml"))
        .toString();
    const mcdrConfig: McdrConfig = jsyaml.load(mcdrConfigYml) as McdrConfig;
    mcdr.mcdrServerPath = mcdrConfig.working_directory;
    mcdr.mcdrPluginsPath = mcdrConfig.plugin_directories;
    this.path = path.join(cwd(), mcdr.mcdrServerPath); // Audit Server's path
    return mcdr;
}

