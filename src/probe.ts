export class Server {
    minecraft: Minecraft;
    fabric?: Fabric;
    forge?: Forge;
    mcdr?: MCDR;

    constructor(
        minecraft: Minecraft,
        fabric?: Fabric,
        forge?: Forge,
        mcdr?: MCDR
    ) {
        this.minecraft = minecraft;
        if (forge && fabric) throw "Coexistence of forge & fabric";
        this.fabric = fabric;
        this.forge = forge;
        this.mcdr = mcdr;
    }
}

class MCDR {
    // serverPath: string;
    // pluginPaths: string[];
    constructor() {}
}

class Minecraft {
    version: string;
    constructor(version: string) {
        this.version = version;
    }
}

class Fabric {
    version: string;
    constructor(version: string) {
        this.version = version;
    }
}

class Forge {
    version: string;
    constructor(version: string) {
        this.version = version;
    }
}

class Bukkit {
    version: string;
    constructor(version: string) {
        this.version = version;
    }
}
