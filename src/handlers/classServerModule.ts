abstract class ServerModule {
    name!: string;
    constructor(name: string) {
        this.name = name;
    }
}

class ServerExecutable {
    path: string;
    minecraft: Minecraft;
    modLoader?: ModLoader;
    constructor(path: string, minecraft: Minecraft, modLoader?: ModLoader) {
        this.path = path;
        this.minecraft = minecraft;
        this.modLoader = modLoader;
    }
}

interface McdrConfig {
    working_directory: string;
    plugin_directories: string[];
}

function isMcdrConfig(obj: unknown): obj is McdrConfig {
    if (typeof obj !== "object" || obj === null) return false;
    const valid =
        typeof (obj as McdrConfig).working_directory === "string" &&
        Array.isArray((obj as McdrConfig).plugin_directories) &&
        (obj as McdrConfig).plugin_directories.every(
            (dir: unknown) => typeof dir === "string"
        );
    return valid;
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

class Minecraft extends ServerModule {
    version: string;
    bukkit?: Bukkit;
    constructor(version: string) {
        super("Minecraft");
        this.version = version;
    }
}

class Fabric extends ServerModule {
    version: string;
    constructor(version: string) {
        super("Fabric");
        this.version = version;
    }
}

class Forge extends ServerModule {
    version: string;
    constructor(version: string) {
        super("Forge");
        this.version = version;
    }
}

class Bukkit extends ServerModule {
    version: string;
    constructor(version: string) {
        super("Bukkit");
        this.version = version;
    }
}

type ModLoader = Fabric | Forge;

export { Mcdr, Minecraft, Fabric, Forge };
export type { McdrConfig, ModLoader };
export { isMcdrConfig };
export { ServerExecutable };
