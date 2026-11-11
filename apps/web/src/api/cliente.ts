import createClient from "openapi-fetch";

import type { paths } from "./esquema";

const baseUrl = import.meta.env.VITE_API_URL?.trim();

export const cliente = createClient<paths>(baseUrl ? { baseUrl } : {});
