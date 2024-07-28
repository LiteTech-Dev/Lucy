import axios from "axios";
import {
    ModrinthVersion,
    ModrinthSearchProjectResult,
    ModrinthProject,
    modrinthVersionHash,
} from "./modrinthTypes.js";

/**
 * Network handler for Modrinth.
 * Singleton.
 */
export class ModrinthApiHandler {
    private static _instance: ModrinthApiHandler;
    private readonly apiRootUrl = "https://api.modrinth.com/v2";

    private constructor() {}

    static getInstance(): ModrinthApiHandler {
        if (!ModrinthApiHandler._instance) {
            ModrinthApiHandler._instance = new ModrinthApiHandler();
        }
        return ModrinthApiHandler._instance;
    }

    /**
     * @param hash File's hash.
     * @param algorithm Hash algorithm used. "sha1" as default.
     * @returns Promise<[1, ModrinthVersion], [Error code, Error]>
     */
    public async getVersionByHash(
        hash: modrinthVersionHash
    ): Promise<ModrinthVersion> {
        const url = `${this.apiRootUrl}/version_file/${hash.val}?algorithm=${hash.algorithm}`;
        return new Promise((resolve, reject) => {
            try {
                console.log(`GET ${url}`);
                axios.get(url).then((res) => {
                    console.log(`SUCCESSFUL: ${res.status} ${res.statusText}`);
                    resolve(res.data);
                });
            } catch (err) {
                return reject(err);
            }
        });
    }

    /**
     * Search a Modrinth Project with specified name and optional facets.
     * @param keyword Name to search.
     * @param facets Search filters. See "Search projects / Query Parameters / facets" in https://docs.modrinth.com/#tag/projects
     * @returns Promise<[1, ModrinthSearchProjectResult], [Error code, Error]>
     */
    public async searchProjects(
        keyword: string,
        facets?: string
    ): Promise<[ModrinthSearchProjectResult]> {
        const url = `${this.apiRootUrl}/search?query=${keyword}${
            facets ? `&facets=${facets}` : ""
        }`;
        return new Promise((resolve, reject) => {
            try {
                console.log(`GET ${url}`);
                axios.get(url).then((res) => {
                    console.log(`SUCCESSFUL: ${res.status} ${res.statusText}`);
                    resolve(res.data);
                });
            } catch (err) {
                return reject(err);
            }
        });
    }

    /**
     * Get a Modrinth project by its ID or slug.
     * @param identifier "project_id" or "slug"
     * @returns Promise<[1, ModrinthProject], [Error code, Error]>
     */
    public async getProject(identifier: string): Promise<[1, ModrinthProject]> {
        const url = `${this.apiRootUrl}/project/${identifier}`;
        return new Promise((resolve, reject) => {
            try {
                console.log(`GET ${url}`);
                axios.get(url).then((res) => {
                    console.log(`SUCCESSFUL: ${res.status} ${res.statusText}`);
                    resolve(res.data);
                });
            } catch (err) {
                return reject(err);
            }
        });
    }
}
