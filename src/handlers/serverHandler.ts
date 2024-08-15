import { cwd } from "process";
import path from "path";
import fs from "fs";
import jsyaml from "js-yaml";
import properties from "properties-reader";

import { ServerFileService } from "../fileUtils/serverFileUtil.js";
import { ServerModule, Mcdr, McdrConfig } from "./classServerModule.js";
import {
    MultipleExecutablesFoundError,
    NoExecutableFoundError,
} from "./serverError.js";

export class ServerHandler {
    private static _instance: ServerHandler;
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
    public static getInstance(): ServerHandler {
        if (!ServerHandler._instance) {
            ServerHandler._instance = new ServerHandler();
        }
        return ServerHandler._instance;
    }

    public moduleExists(moduleName: string): boolean {
        return this.modules.some((module) => module.name === moduleName);
    }

    private getModule(moduleName: string): ServerModule | null {
        return (
            this.modules.find((module) => module.name === moduleName) || null
        );
    }

    private searchForServerExecutable(): string {
        const suspectedServerExecutable: string[] = [];
        const validExecutables: string[] = [];

        // 获取当前路径下的所有.jar文件
        const files = fs.readdirSync(this.path);
        files.forEach((file) => {
            if (path.extname(file) === ".jar") {
                suspectedServerExecutable.push(path.join(this.path, file));
            }
        });

        // 拆包检查
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
            // TODO: 检测forge jar
        });

        if (validExecutables.length > 1)
            throw new MultipleExecutablesFoundError();
        if (validExecutables.length === 0) throw new NoExecutableFoundError();

        // 返回找到的服务器可执行文件的路径
        return validExecutables[0];
    }

    private readMcdrConfig(): McdrConfig {
        const mcdrConfigPath = path.resolve(cwd(), "config.yml");
        if (!fs.existsSync(mcdrConfigPath))
            throw new Error("MCDR config file not found");
        const mcdrConfigYml: string = fs
            .readFileSync(mcdrConfigPath)
            .toString();
        const mcdrConfig = jsyaml.load(mcdrConfigYml) as McdrConfig;
        if (Object.values(mcdrConfig).some((value) => !value))
            throw new Error("Invalid MCDR config file");
        return mcdrConfig;
    }
}
