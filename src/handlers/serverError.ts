abstract class ServerError extends Error {
    constructor(message: string) {
        super(message);
        this.name = "ServerError";
    }
}

abstract class ServerExecutableError extends ServerError {
    constructor(message: string) {
        super(message);
        this.name = "ServerExecutableError";
    }
}

class NoExecutableFoundError extends ServerExecutableError {
    constructor() {
        super("No server executable found");
        this.name = "NoExecutableFoundError";
    }
}

class MultipleExecutablesFoundError extends ServerExecutableError {
    constructor() {
        super("Multiple server executables found");
        this.name = "MultipleExecutablesFoundError";
    }
}

abstract class McdrError extends ServerError {
    constructor(message: string) {
        super(message);
        this.name = "McdrError";
    }
}

class McdrConfigNotFoundError extends McdrError {
    constructor() {
        super("MCDR config not found");
        this.name = "McdrConfigNotFoundError";
    }
}

class McdrConfigInvalidError extends McdrError {
    constructor() {
        super("MCDR config is invalid");
        this.name = "McdrConfigInvalidError";
    }
}

abstract class ModLoaderError extends ServerError {
    constructor(message: string) {
        super(message);
        this.name = "ModLoaderError";
    }
}

class CoexistenceOfFabricAndForgeError extends ModLoaderError {
    constructor() {
        super("There are multiple mod loaders found");
        this.name = "CoexistenceOfFabricAndForgeError";
    }
}

abstract class MinecraftError extends ServerError {
    constructor(message: string) {
        super(message);
        this.name = "MinecraftError";
    }
}

class MinecraftNotFoundError extends MinecraftError {
    constructor() {
        super("Minecraft not found");
        this.name = "MinecraftNotFoundError";
    }
}

export {
    NoExecutableFoundError,
    MultipleExecutablesFoundError,
    McdrConfigNotFoundError,
    McdrConfigInvalidError,
    CoexistenceOfFabricAndForgeError,
    MinecraftNotFoundError,
};
