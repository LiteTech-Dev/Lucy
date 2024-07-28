export type modrinthHashAlgorithm = "sha1" | "sha512";
export type modrinthVersionHash = {
    val: string;
    algorithm: modrinthHashAlgorithm;
};

export type ModrinthProject = {
    slug: string[]; // my_project
    title: string; // My project
    description: string; // A short description
    categories: string[];
    client_side: "required" | "optional" | "unsupported";
    server_side: "required" | "optional" | "unsupported";
    body: string; // A long body describing my project in detail
    status: string;
    requested_status: string;
    additional_categories: string[];
    issues_url: string;
    source_url: string;
    wiki_url: string;
    discord_url: string;
    donation_urls?: { id: string; platform: string; url: string }[];
    project_type: string;
    downloads: number;
    icon_url: string;
    color: number;
    thread_id: string;
    monetization_status?: string;
    id: string;
    team: string;
    body_url?: string | null;
    moderator_message?: string | null;
    published: string;
    updated: string;
    approved: string;
    queued: string;
    followers: number;
    lisense?: { id: string; name: string; url?: string };
    versions: string[]; // Version hash
    game_versions: string[];
    loaders: (
        | "forge"
        | "fabric"
        | "quilt"
        | "neoforge"
        | "lightloader"
        | "rift"
    )[];
    gallery: {
        url: string;
        featured: boolean;
        title: string;
        description: string;
        created: string;
        ordering: number;
    }[];
};

// /**
//  * List of valid facets for searching Modrinth project.
//  */
// enum ModrinthSearchValidFacets {
//     project_type = "project_type", // string[]
//     categories = "categories", // string[] (including "forge" and "fabric" etc.)
//     versions = "versions", // string[]
//     client_side = "client_side", // "required" | "optional" | "unsupported"
//     server_side = "server_side", // "required" | "optional" | "unsupported"
//     open_source = "open_source", // boolean
//     title = "title", // string
//     author = "author", // string
//     follows = "follows", // number
//     project_id = "project_id", // string
//     license = "license", // string
//     downloads = "downloads", // number
//     color = "color", // number
//     created_timestamp = "created_timestamp",
//     modified_timestamp = "modified_timestamp",
// }

// /**
//  * Template of Modrinth searching facets argument.
//  */
// type ModrinthSearchArgs = {
//     argName: ModrinthSearchValidFacets;
//     argValue: (string | string[] | boolean | number)[];
// };

export type ModrinthSearchHits = {
    project_id: string;
    project_type: string;
    slug: string;
    author: string;
    title: string;
    description: string;
    categories: string[];
    display_categories: string[];
    versions: string[]; // Version number such as "1.12.2"
    downloads: number;
    follows: number;
    icon_url: string;
    date_created: string; // Timestamp like "2023-01-29T00:32:51.714996Z"
    date_modified: string; // Timestamp
    latest_version: string; // Version ID I think
    license: string;
    client_side: "required" | "optional" | "unsupported";
    server_side: "required" | "optional" | "unsupported";
    gallery: string[]; // URLs
    featured_gallery: null | unknown; // idk, all samples are "null"
    color: number;
};

export type ModrinthSearchProjectResult = {
    hits: ModrinthSearchHits[];
    offset: number;
    limit: number;
    total_hits: number;
};

/**
 * Raw format of JSON that Modrinth "version" api responses.
 * You should convert it to class.mod
 */
export type ModrinthVersion = {
    game_versions: string[];
    loaders: string[];
    id: string;
    project_id: string;
    featured: boolean;
    name: string;
    version_number: string;
    changelog?: string;
    changelog_url?: string;
    date_published: string;
    downloads: number;
    version_type: "alpha" | "beta" | "release";
    status: string;
    requested_status: string;
    files: object[];
    dependencies: string[];
};
