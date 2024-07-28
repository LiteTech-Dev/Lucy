import axios from "axios"
/**
 * Network handler for Modrinth.
 * Singleon, use ModrinthHandler.getInstance to get an instance.
 */
export class ModrinthHandler{
    private static _instance: ModrinthHandler

    private constructor(){} //Singleton

    static getInstance(): ModrinthHandler{
        if(!ModrinthHandler._instance){
            ModrinthHandler._instance = new ModrinthHandler()
            console.log("Created ModrinthHandler instance")
        }
        return ModrinthHandler._instance
    }

    /**
     * Use file's hash to get Modrinth version object.
     * try...catch is required when using this function.
     * @param hash Hash of the file.
     * @param algorithm Hash algorithm used. "sha1" as default.
     */
    public async searchVersionHash(hash: string, algorithm:"sha1"|"sha512" ="sha1"): Promise<[boolean, object | unknown]>{
        const fullUrl= `https://api.modrinth.com/v2/version_file/${hash}?algorithm=${algorithm}`
        console.log(`GET ${fullUrl}`);
        try{
            const response:[true, object] = [ true, JSON.parse((await axios.get(fullUrl)).data) ]
            return response
        }
        catch(err){
            const response:[false, unknown] = [ false, err ]
            return response
        }


    }
}

