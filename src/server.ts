class Server {
    modules: { [key: string]: ServerModule } = {};
}

interface ServerModule {
    version?: string;
}

export { Server };
export type { ServerModule };
