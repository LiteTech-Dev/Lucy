import { cwd } from "process";
import path from "path";
import fs from "fs";
import jsyaml from "js-yaml";
// import properties from "properties-reader";

import { ServerFileService } from "../fileUtils/serverFileUtil.js";
import {
    isMcdrConfig,
    Mcdr,
    McdrConfig,
    Minecraft,
    ModLoader,
} from "./classServerModule.js";
import {
    McdrConfigInvalidError,
    McdrConfigNotFoundError,
    MultipleExecutablesFoundError,
    NoExecutableFoundError,
} from "./serverError.js";

export class ServerHandler {
    private static _instance: ServerHandler;
    private modLoader?: ModLoader;
    private mcdr?: Mcdr;
    private minecraft: Minecraft = new Minecraft("0");
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
            // TODO: 设计检测逻辑，是否顺便检测ModLoader
            this.executable = this.searchForServerExecutable();
        } catch (err) {
            // TODO: Return this to the user
            // If no executable is found, exit and notify the user
            // If multiple executables are found, let them choose a default
        }
    }

    // Global access point
    public static getInstance(): ServerHandler {
        if (!ServerHandler._instance) {
            ServerHandler._instance = new ServerHandler();
        }
        return ServerHandler._instance;
    }

    private searchForServerExecutable(): string {
        const suspectedServerExecutable: string[] = [];
        const validExecutables: string[] = [];

        // Add all .jar files to suspectedServerExecutable[]
        const files = fs.readdirSync(this.path);
        files.forEach((file) => {
            if (path.extname(file) === ".jar") {
                suspectedServerExecutable.push(path.join(this.path, file));
            }
        });

        // 拆包检查
        // TODO: 编写检查逻辑
        suspectedServerExecutable.forEach((file) => {
            // 检测原版jar
            ServerFileService.getInstance().readFileFromZip(
                file,
                "version.json"
            );
            // 检测fabric jar
            ServerFileService.getInstance().readFileFromZip(
                file,
                "install.properties"
            );
            // TODO: 读Forge文档
        });

        switch (validExecutables.length) {
            case 0:
                throw new NoExecutableFoundError();
            case 1:
                return validExecutables[0];
            default:
                throw new MultipleExecutablesFoundError();
        }
    }
}

export function readMcdrConfig(): McdrConfig {
    const mcdrConfigPath = path.resolve(cwd(), "config.yml");
    if (!fs.existsSync(mcdrConfigPath)) throw new McdrConfigNotFoundError();
    const mcdrConfigObject = jsyaml.load(
        fs.readFileSync(mcdrConfigPath, "utf-8").toString()
    );

    if (!isMcdrConfig(mcdrConfigObject)) throw new McdrConfigInvalidError();

    // Sanitize the config object
    for (const key of Object.keys(mcdrConfigObject) as Array<
        keyof typeof mcdrConfigObject
    >) {
        if (!["working_directory", "plugin_directories"].includes(key)) {
            delete mcdrConfigObject[key];
        }
    }

    return mcdrConfigObject;
}
