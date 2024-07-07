import path from "path";
import os from "os";
import { cwd } from "process";

// A .lucy path exists for all server that was initialized by Lucy.
// It works similar to .git.
// All local cache, configuration, and other files are stored in this directory.
export const LUCY_PATH = path.join(cwd(), ".lucy");

// The global cache path
// For network requests, file downloads, etc.
export const CACHE_PATH = path.join(os.tmpdir(), ".lucy_cache");
