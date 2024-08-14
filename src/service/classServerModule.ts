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

class Minecraft extends ServerModule {
    version: string;
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

export { ServerModule, Mcdr, Minecraft, Fabric, Forge };
export type { McdrConfig };
