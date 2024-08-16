import fs from "fs";
import {
    McdrConfigNotFoundError,
    McdrConfigInvalidError,
} from "../serverError";
import { readMcdrConfig } from "../serverHandler.js";

describe("ServerHandler", () => {
    describe("readMcdrConfig", () => {
        beforeEach(() => {
            jest.spyOn(process, "cwd").mockReturnValue(__dirname);
        });

        afterEach(() => {
            jest.restoreAllMocks();
        });

        it("should read the MCDR config file and return the parsed config object", () => {
            jest.spyOn(fs, "existsSync").mockReturnValue(true);

            expect(readMcdrConfig()).toEqual({
                working_directory: "server",
                plugin_directories: ["plugins"],
            });
        });

        it("should throw McdrConfigNotFoundError if the MCDR config file does not exist", () => {
            jest.spyOn(fs, "existsSync").mockReturnValue(false);

            expect(() => readMcdrConfig()).toThrow(
                new McdrConfigNotFoundError()
            );
        });

        it("should throw McdrConfigInvalidError if the MCDR config file is invalid", () => {
            jest.spyOn(fs, "existsSync").mockReturnValue(true);
            jest.spyOn(fs, "readFileSync").mockReturnValue(`
                invalid_config: 123
                working_directory: 
                plugin_directories: 
            `);

            expect(() => readMcdrConfig()).toThrow(
                new McdrConfigInvalidError()
            );
        });
    });
});
