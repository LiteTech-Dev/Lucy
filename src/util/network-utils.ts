import axios, { AxiosError, AxiosResponse } from "axios";
import {
    ModrinthVersion,
    ModrinthSearchProjectResult,
    ModrinthProject,
} from "../class/modrinth-objects.js";

/**
 *
 * @param err Input the catched error.
 * @returns HTTP status code | 2 -> No response from remote server | 3 -> Unknown network error | 4 -> Error when sending request (not a network error)
 */
function getNetWorkErrorCode(err: AxiosError | unknown) {
    let _errCode: number;
    if (err instanceof AxiosError) {
        if (err.response) {
            console.log(
                "FAILED: " +
                    err.response.status.toString() +
                    " " +
                    err.response.statusText.toString()
            );
            _errCode = err.response.status;
        } else if (err.request) {
            console.log("FAILED: No response from remote server");
            _errCode = 2;
        } else {
            console.log(
                "FAILED: Unknown network error \n" + err.toJSON().toString()
            );
            _errCode = 3;
        }
    } else {
        console.log("FAILED: Error when sending request: ");
        console.log(err);
        _errCode = -1;
    }
    return _errCode;
}

// async function requestGET(fullUrl: string) {
//     try {
//         const response: AxiosResponse = await axios.get(fullUrl);
//         console.log(
//             "SUCCESSFUL: " +
//                 response.status.toString() +
//                 response.statusText.toString()
//         );
//         return [1, response.data];
//     } catch (err: unknown | AxiosError) {
//         return [getNetWorkErrorCode(err), err];
//     }
// }

/**
 * Network handler for Modrinth.
 * Singleton.
 */
export class ModrinthApiHandler {
    private static _instance: ModrinthApiHandler;

    private apiRootUrl = "https://api.modrinth.com/v2";
    //private apiRootUrlDev = "https://stage.api.modrinth.com/v2"

    private constructor() {} //Singleton

    static getInstance(): ModrinthApiHandler {
        if (!ModrinthApiHandler._instance) {
            ModrinthApiHandler._instance = new ModrinthApiHandler();
        }
        return ModrinthApiHandler._instance;
    }

    /**
     *
     * @param hash File's hash.
     * @param algorithm Hash algorithm used. "sha1" as default.
     * @returns Promise<[1, class/modrinth-objects/ModrinthVersion], [Error code, Error]>
     */
    public async getVersionByHash(
        hash: string,
        algorithm: "sha1" | "sha512" = "sha1"
    ): Promise<[1, ModrinthVersion]> {
        const _fullUrl: string =
            this.apiRootUrl + `/version_file/${hash}?algorithm=${algorithm}`;

        try {
            console.log("GET " + _fullUrl);
            const response: AxiosResponse = await axios.get(_fullUrl);
            console.log(
                "SUCCESSFUL: " +
                    response.status.toString() +
                    " " +
                    response.statusText.toString()
            );
            return [1, response.data];
        } catch (err: unknown | AxiosError) {
            return Promise.reject([getNetWorkErrorCode(err), err]);
        }
    }

    /**
     * Search a Modrinth Project with specified name and optional facets.
     * @param query Name to search.
     * @param facets Search filters. See "Search projects / Query Parameters / facets" in https://docs.modrinth.com/#tag/projects
     * @returns Promise<[1, class/modrinth-objects/ModrinthSearchProjectResult], [Error code, Error]>
     */
    public async searchProjects(
        query: string,
        facets?: string
    ): Promise<[1, ModrinthSearchProjectResult]> {
        const _fullUrl: string = this.apiRootUrl + `/search?query=${query}`;
        let fullUrl = _fullUrl;

        if (facets) {
            fullUrl = fullUrl + "&facets=" + `${facets}`;
        }

        try {
            console.log("GET " + fullUrl);
            const response: AxiosResponse = await axios.get(fullUrl);
            console.log(
                "SUCCESSFUL: " +
                    response.status.toString() +
                    " " +
                    response.statusText.toString()
            );
            return [1, response.data];
        } catch (err: unknown | AxiosError) {
            return Promise.reject([getNetWorkErrorCode(err), err]);
        }
    }

    /**
     * Get a Modrinth project by its ID or slug.
     * @param identifier "project_id" or "slug"
     * @returns Promise<[1, class/modrinth-objects/ModrinthProject], [Error code, Error]>
     */
    public async getProject(identifier: string): Promise<[1, ModrinthProject]> {
        const _fullUrl: string = this.apiRootUrl + `/project/${identifier}`;
        try {
            console.log("GET " + _fullUrl);
            const response: AxiosResponse = await axios.get(_fullUrl);
            console.log(
                "SUCCESSFUL: " +
                    response.status.toString() +
                    " " +
                    response.statusText.toString()
            );
            return [1, response.data];
        } catch (err: unknown | AxiosError) {
            return Promise.reject([getNetWorkErrorCode(err), err]);
        }
    }
}
