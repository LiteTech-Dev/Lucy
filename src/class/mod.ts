class Mod {
    public properties: {
        name?: string;
        fileName?: string;
        filePath?: string; //TODO: Use path.join
        version?: string;
        author?: string;
        hash?: string;
        loader?: string;
        gameVersions?: string[];
        dependencies?: Mod[];
        versionType?: string;
        changelog?: string;
        changelogUrl?: string;
        isCompatible?: boolean;
    } = {};

    // public function getModrinthVersion(args) {
    //     new ModrinthApiCompatible
    //     return ModrinthApiCompatible
    // }

    // public function getCurseForgeCompatible(params:type) {
    //     new CurseForgeCompatible
    //     return CurseForgeCompatible
    // }
}

export { Mod };
