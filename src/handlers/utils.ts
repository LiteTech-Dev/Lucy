import jsyaml from "js-yaml";
import path from "path";
import fs from "fs";
import { cwd } from "process";
import { isMcdrConfig, McdrConfig } from "./classServerModule.js";
import {
    McdrConfigNotFoundError,
    McdrConfigInvalidError,
    MultipleExecutablesFoundError,
    NoExecutableFoundError,
} from "./serverError.js";
import { ServerFileService } from "../fileUtils/serverFileUtil.js";

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

export function searchForServerExecutable(pathToSearch: string): string {
    const suspectedServerExecutable: string[] = [];
    const validExecutables: string[] = [];

    // Add all .jar files to suspectedServerExecutable[]
    const files = fs.readdirSync(pathToSearch);
    files.forEach((file) => {
        if (path.extname(file) === ".jar") {
            suspectedServerExecutable.push(path.join(pathToSearch, file));
        }
    });

    // 拆包检查
    // TODO: 编写检查逻辑
    suspectedServerExecutable.forEach((file) => {
        // 检测原版jar
        ServerFileService.getInstance().readFileFromZip(file, "version.json");
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
