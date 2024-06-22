class ServerSha1 {
    version: string;
    sha1: string;
    constructor(version: string, sha1: string) {
        this.sha1 = sha1;
        this.version = version;
    }
}

async function updateSha1Cache(): Promise<void> {
    // check the mojang api, if the newest version changed, update
    // or, if cache does not exist, create one
    // else, return
}

async function createSha1Cache(): Promise<void> {
    const manifestData: any = await fetch(
        "https://launchermeta.mojang.com/mc/game/"
    );

    let servers: ServerSha1[];

    manifestData.versions.forEach(async (version: any) => {
        const versionData: any = await fetch(version.url);
        if (versionData.downloads.server.sha1) {
            servers.push(
                new ServerSha1(version.id, versionData.downloads.server.sha1)
            );
        }
    });
}
