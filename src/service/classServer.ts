import { cwd } from "process";
import path from "path";
import jsyaml from "js-yaml";
import fs from "fs";

export class Server {
    private static _instance: Server;
    path: string;
    private modules: ServerModule[] = [];

    private constructor() {
        // Initialize MCDR
        const mcdrConfig: McdrConfig | null = this.readMcdrConfig();

        if (mcdrConfig) {
            const mcdr = new Mcdr(mcdrConfig);
            this.modules.push(mcdr);
            this.path = mcdr.mcdrServerPath;
        } else {
            this.path = cwd();
        }
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
    private readMcdrConfig(): McdrConfig | null {
        const mcdrConfigPath = path.resolve(cwd(), "config.yml");
        let mcdrConfig: McdrConfig;

        try {
            const mcdrConfigYml: string = fs
                .readFileSync(mcdrConfigPath)
                .toString();
            mcdrConfig = jsyaml.load(mcdrConfigYml) as McdrConfig;
            if (Object.values(mcdrConfig).some((value) => !value)) return null;
            return mcdrConfig;
        } catch (error) {
            // TODO: Log error properly
            console.error("Error reading config:", error);
        }

        // Return null if reading config fails
        return null;
    }
}

class ServerModule {
    name!: string;
    constructor(name: string) {
        this.name = name;
    }
}

interface McdrConfig {
    working_directory: string;
    plugin_directories: string[];
}

class Mcdr extends ServerModule {
    mcdrServerPath: string = "";
    mcdrPluginPaths: string[] = [];

    constructor(config: McdrConfig) {
        super("MCDR");
        this.mcdrServerPath = config.working_directory;
        this.mcdrPluginPaths = config.plugin_directories;
    }
}
}

