import * as fs from "fs";
import * as yaml from "js-yaml";

interface MCDR {
    // version: string;
    // I don't think we should or are able to check MCDR itselves version.
    serverPath: string;
    pluginPaths: string[];
}

export function checkMCDR(path: string = "."): MCDR | null {
    const mcdrConfigFilePath = path + "/config.yml";
    if (fs.existsSync(mcdrConfigFilePath)) {
        const mcdrConfigYAML = fs.readFileSync(mcdrConfigFilePath, "utf-8");
        const mcdrConfig = yaml.load(mcdrConfigYAML) as {
            working_directory: string;
            plugin_directories: string[];
        };
        const MCDR: MCDR = {
            serverPath: mcdrConfig.working_directory,
            pluginPaths: mcdrConfig.plugin_directories,
        };
        return MCDR;
    }
    return null;
}

// interface Minecraft {
//     version: string;
// }

// export function checkMinecraft(): Minecraft | null {}
// Unsure how to implement yet.

// interface Fabric {
//     version: string;
// }

// export function checkFabric(path: string = "."): Fabric | null {
//     const fabricLoaderPath = "libraries/net/fabricmc/fabric-loader";
// }
