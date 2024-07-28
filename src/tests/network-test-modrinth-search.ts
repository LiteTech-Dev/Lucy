import { ModrinthApiHandler } from "../networkUtils/modrinthApiHandler.js";
export async function runTest() {
    ModrinthApiHandler.getInstance()
        .searchProjects("Better", '[["author:LingLing1301"]]')
        .then(function (result) {
            console.log(result);
        })
        .catch(function (err) {
            console.error(err);
        });
}

runTest();
