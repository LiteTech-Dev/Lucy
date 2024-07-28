class Server {
    modules: { [key: string]: ServerModule } = {};
}

type ServerModule = {
    version?: string;
};

export { Server };
export type { ServerModule };
