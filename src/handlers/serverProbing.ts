import jsyaml from "js-yaml";
import path from "path";
import fs from "fs";
import { cwd } from "process";
import {
    Fabric,
    isMcdrConfig,
    McdrConfig,
    Minecraft,
    ServerExecutable,
} from "./classServerModule.js";
import {
    McdrConfigNotFoundError,
    McdrConfigInvalidError,
    MultipleExecutablesFoundError,
    NoExecutableFoundError,
    InvalidExecutableError,
} from "./serverError.js";
import { ZipUtil } from "../fileUtils/zipUtil.js";

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

function searchForSuspectedServerExecutable(pathToSearch: string): string[] {
    const suspectedServerExecutable: string[] = [];
    const files = fs.readdirSync(pathToSearch);
    files.forEach((file) => {
        if (path.extname(file) === ".jar") {
            suspectedServerExecutable.push(path.join(pathToSearch, file));
        }
    });
    return suspectedServerExecutable;
}

async function validiateServerExecutable(
    executablePath: string
): Promise<ServerExecutable> {
    if (
        await ZipUtil.getInstance().fileExistsInZip(
            executablePath,
            "version.json"
        )
    ) {
        const data = await ZipUtil.getInstance().readFileFromZip(
            executablePath,
            "version.json"
        );
        return new ServerExecutable(executablePath, initVanilla(data));
    } else if (
        await ZipUtil.getInstance().fileExistsInZip(
            executablePath,
            "install.properties"
        )
    ) {
        const data = await ZipUtil.getInstance().readFileFromZip(
            executablePath,
            "install.properties"
        );
        const [minecraft, fabric] = initFabric(data);
        return new ServerExecutable(executablePath, minecraft, fabric);
    }
    throw new InvalidExecutableError();
}

export function searchForServerExecutable(
    pathToSearch: string
): ServerExecutable {
    // This function was still poorly defined
    // The process of validiating the executable includes determining the version and other imformation

    // Now the solution was defining a new data structure ServerExecutable, and returning it as a whole

    const suspectedServerExecutables: string[] =
        searchForSuspectedServerExecutable(pathToSearch);
    const validExecutables: ServerExecutable[] = [];

    suspectedServerExecutables.forEach((executable) => {
        try {
            validiateServerExecutable(executable).then((serverExecutable) => {
                validExecutables.push(serverExecutable);
            });
        } catch (err) {
            if (err instanceof InvalidExecutableError) {
                // Do nothing
            } else {
                throw err;
            }
        }
    });

    switch (validExecutables.length) {
        case 0:
            throw new NoExecutableFoundError();
        case 1:
            return validExecutables[0];
        default:
            throw new MultipleExecutablesFoundError();
        // After throwing this, the caller should prompt the user to choose a default
        // Then, the user's choice should be saved in the config file
    }
}

function initVanilla(data: string): Minecraft {
    return new Minecraft(JSON.parse(data).id);
}

function initFabric(data: string): [Minecraft, Fabric] {
    const properties: string[] = data.split("\n");
    const fabricLoaderVersion: string = properties[0].split("=")[1];
    const gameVersion: string = properties[1].split("=")[1];
    return [new Minecraft(gameVersion), new Fabric(fabricLoaderVersion)];
}
