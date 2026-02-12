export interface JRPC {
    login(): Promise<number>
    create(model: string, data: any): Promise<any>
    load(model: string, header: any, values: any): Promise<any>
    count(model: string, domain?: string): Promise<number>
    fields_get(model: string, fields?: string[], field_attributes?: string[]): Promise<any>
    get_id(model: string, domain?: string): Promise<number>
    search(model: string, domain?: string): Promise<number[]>
    read<T>(model: string, ids: number[], fields?: string[]): Promise<T[]>
    search_read<T>(model: string, domain?: string, fields?: string[], offset?: number, limit?: number): Promise<T[]>
    write(model: string, id: number | number[], data: any): Promise<any>
    unlink(model: string, ids: number[]): Promise<any>
    execute(model: string, method: string, args: any): Promise<any>
    execute_kw(model: string, method: string, args?: any[], kwargs?: any): Promise<any>
}

export class OdooRPCError extends Error {
    constructor(
        message: string,
        public code?: number,
        public data?: any
    ) {
        super(message);
        this.name = 'OdooRPCError';
    }
}

export class JSONRPC implements JRPC {
    constructor(private config: JRPCConfig) {
        config.url = `${config.schema}://${config.hostname}:${config.port}/jsonrpc/`;
    }

    async login(): Promise<number> {
        const [uid, error] = await this.Call("common", "login", [
            this.config.database,
            this.config.username,
            this.config.password
        ]);

        if (error) {
            throw new OdooRPCError(error.message || 'Login failed', error.code, error.data);
        }

        if (!uid) {
            throw new OdooRPCError('Login failed: Invalid credentials');
        }

        this.config.uid = uid;
        return uid;
    }

    async create(model: string, data: any): Promise<any> {
        const [result, error] = await this.Call("object", "execute", [
            this.config.database,
            this.config.uid,
            this.config.password,
            model,
            "create",
            data
        ]);

        if (error) {
            throw new OdooRPCError(error.message || 'Create failed', error.code, error.data);
        }

        return result;
    }

    async load(model: string, header: any, values: any): Promise<any> {
        const [result, error] = await this.Call("object", "execute", [
            this.config.database,
            this.config.uid,
            this.config.password,
            model,
            "load",
            [header, values]
        ]);

        if (error) {
            throw new OdooRPCError(error.message || 'Load failed', error.code, error.data);
        }

        return result;
    }

    async count(model: string, domain?: string): Promise<number> {
        const [result, error] = await this.Call("object", "execute", [
            this.config.database,
            this.config.uid,
            this.config.password,
            model,
            "search_count",
            parseDomainString(domain ?? "")
        ]);

        if (error) {
            throw new OdooRPCError(error.message || 'Count failed', error.code, error.data);
        }

        return result;
    }

    async fields_get(model: string, fields?: string[], field_attributes?: string[]): Promise<any> {
        const [result, error] = await this.Call("object", "execute", [
            this.config.database,
            this.config.uid,
            this.config.password,
            model,
            "fields_get",
            fields ?? [],
            field_attributes ?? []
        ]);

        if (error) {
            throw new OdooRPCError(error.message || 'fields_get failed', error.code, error.data);
        }

        return result;
    }

    async get_id(model: string, domain?: string): Promise<number> {
        const [result, error] = await this.Call("object", "execute", [
            this.config.database,
            this.config.uid,
            this.config.password,
            model,
            "search",
            parseDomainString(domain ?? "")
        ]);

        if (error) {
            throw new OdooRPCError(error.message || 'GetID failed', error.code, error.data);
        }

        return result[0] || -1;
    }

    async search(model: string, domain?: string): Promise<number[]> {
        const [result, error] = await this.Call("object", "execute", [
            this.config.database,
            this.config.uid,
            this.config.password,
            model,
            "search",
            parseDomainString(domain ?? "")
        ]);

        if (error) {
            throw new OdooRPCError(error.message || 'Search failed', error.code, error.data);
        }

        return result;
    }

    async read<T>(model: string, ids: number[], fields?: string[]): Promise<T[]> {
        const [result, error] = await this.Call("object", "execute", [
            this.config.database,
            this.config.uid,
            this.config.password,
            model,
            "read",
            ids,
            fields ?? []
        ]);

        if (error) {
            throw new OdooRPCError(error.message || 'Read failed', error.code, error.data);
        }

        return result;
    }

    async search_read<T>(
        model: string,
        domain?: string,
        fields?: string[],
        offset?: number,
        limit?: number
    ): Promise<T[]> {
        const [result, error] = await this.Call("object", "execute", [
            this.config.database,
            this.config.uid,
            this.config.password,
            model,
            "search_read",
            parseDomainString(domain ?? ""),
            fields ?? [],
            offset ?? 0,
            limit ?? 0
        ]);

        if (error) {
            throw new OdooRPCError(error.message || 'SearchRead failed', error.code, error.data);
        }

        return result;
    }

    async write(model: string, id: number | number[], data: any): Promise<any> {
        // Ensure ids are in array format
        const ids = Array.isArray(id) ? id : [id];

        const [result, error] = await this.Call("object", "execute", [
            this.config.database,
            this.config.uid,
            this.config.password,
            model,
            "write",
            ids,
            data
        ]);

        if (error) {
            throw new OdooRPCError(error.message || 'Write failed', error.code, error.data);
        }

        return result;
    }

    async unlink(model: string, ids: number[]): Promise<any> {
        const [result, error] = await this.Call("object", "execute", [
            this.config.database,
            this.config.uid,
            this.config.password,
            model,
            "unlink",
            ids
        ]);

        if (error) {
            throw new OdooRPCError(error.message || 'Unlink failed', error.code, error.data);
        }

        return result;
    }

    async execute(model: string, method: string, args: any): Promise<any> {
        const [result, error] = await this.Call("object", "execute", [
            this.config.database,
            this.config.uid,
            this.config.password,
            model,
            method,
            args
        ]);

        if (error) {
            throw new OdooRPCError(error.message || `Execute ${method} failed`, error.code, error.data);
        }

        return result;
    }

    async execute_kw(model: string, method: string, args: any[] = [], kwargs: any = {}): Promise<any> {
        const [result, error] = await this.Call("object", "execute_kw", [
            this.config.database,
            this.config.uid,
            this.config.password,
            model,
            method,
            args,
            kwargs
        ]);

        if (error) {
            throw new OdooRPCError(error.message || `ExecuteKw ${method} failed`, error.code, error.data);
        }

        return result;
    }

    // 🔹 Protocol-specific extras
    async Call(service: string, method: string, args: any[] = []): Promise<[any, any]> {
        const req = this.encodeClientRequest(service, method, args);

        try {
            const response = await fetch(this.config.url!, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify(req),
            });

            if (!response.ok) {
                throw new OdooRPCError(
                    `HTTP Error: ${response.status} ${response.statusText}`,
                    response.status
                );
            }

            const json = await response.json() as { result: any; error: any };
            return [json.result, json.error];
        } catch (error) {
            if (error instanceof OdooRPCError) {
                throw error;
            }
            throw new OdooRPCError(
                `Network error: ${error instanceof Error ? error.message : 'Unknown error'}`
            );
        }
    }

    encodeClientRequest(service: string, method: string, args: any[] = []): any {
        return {
            "jsonrpc": "2.0",
            "method": "call",
            "id": Math.floor(Math.random() * 1_000_000_000),
            "params": {
                "service": service,
                "method": method,
                "args": [...args],
            },
        };
    }
}

export class ODOOJSON implements JRPC {
    constructor(private config: JRPCConfig) {
        config.url = `${config.schema}://${config.hostname}:${config.port}/json/2/`;
    }

    async login(): Promise<number> {
        // JSON2RPC protocol does not require a login call
        // Authentication is handled via API key in headers
        // Return -1 to indicate "no UID needed" (not a real user)
        // Do NOT return a real UID like 1 (odoobot) or 2 (admin)
        return -1;
    }

    async create(model: string, data: any): Promise<any> {
        const url = `${this.config.url}${model}/create`;

        try {
            const response = await fetch(url, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    "Authorization": `Bearer ${this.config.apikey}`,
                },
                body: JSON.stringify({ data: data }),
            });

            if (!response.ok) throw new Error(`HTTP error: ${response.status}`);

            const result = await response.json();
            return result;
        } catch (error) {
            throw new OdooRPCError(
                `Create failed: ${error instanceof Error ? error.message : 'Unknown error'}`
            );
        }
    }

    async load(model: string, header: any, values: any): Promise<any> {
        const url = `${this.config.url}${model}/load`;

        try {
            const response = await fetch(url, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    "Authorization": `Bearer ${this.config.apikey}`,
                },
                body: JSON.stringify({ header: header, values: values }),
            });

            if (!response.ok) throw new Error(`HTTP error: ${response.status}`);

            const result = await response.json();
            return result;
        } catch (error) {
            throw new OdooRPCError(
                `Load failed: ${error instanceof Error ? error.message : 'Unknown error'}`
            );
        }
    }

    async count(model: string, domain?: string): Promise<number> {
        const url = `${this.config.url}${model}/count`;

        try {
            const response = await fetch(url, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    "Authorization": `Bearer ${this.config.apikey}`,
                },
                body: JSON.stringify({ domain: parseDomainString(domain ?? "") }),
            });

            if (!response.ok) throw new Error(`HTTP error: ${response.status}`);

            const data = await response.json() as { count: number };
            return data.count;
        } catch (error) {
            throw new OdooRPCError(
                `Count failed: ${error instanceof Error ? error.message : 'Unknown error'}`
            );
        }
    }

    async fields_get(model: string, fields?: string[], field_attributes?: string[]): Promise<any> {
        const url = `${this.config.url}${model}/fields_get`;

        try {
            const response = await fetch(url, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    "Authorization": `Bearer ${this.config.apikey}`,
                },
                body: JSON.stringify({
                    context: {},
                    allfields: fields ? fields : undefined,
                    attributes: field_attributes ? field_attributes : undefined,
                }),
            });

            if (!response.ok) throw new Error(`HTTP error: ${response.status}`);

            const data = await response.json();
            return data;
        } catch (error) {
            throw new OdooRPCError(
                `Fields_get failed: ${error instanceof Error ? error.message : 'Unknown error'}`
            );
        }
    }

    async get_id(model: string, domain?: string): Promise<number> {
        const url = `${this.config.url}${model}/search`;

        try {
            const response = await fetch(url, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    "Authorization": `Bearer ${this.config.apikey}`,
                },
                body: JSON.stringify({ domain: parseDomainString(domain ?? "") }),
            });

            if (!response.ok) throw new Error(`HTTP error: ${response.status}`);

            const data = await response.json() as { ids: number[] };
            return data.ids[0] ?? -1;
        } catch (error) {
            throw new OdooRPCError(
                `Get_id failed: ${error instanceof Error ? error.message : 'Unknown error'}`
            );
        }
    }

    async search(model: string, domain?: string): Promise<number[]> {
        const url = `${this.config.url}${model}/search`;

        try {
            const response = await fetch(url, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    "Authorization": `Bearer ${this.config.apikey}`,
                },
                body: JSON.stringify({ domain: parseDomainString(domain ?? "") }),
            });

            if (!response.ok) throw new Error(`HTTP error: ${response.status}`);

            const data = await response.json() as { ids: number[] };
            return data.ids;
        } catch (error) {
            throw new OdooRPCError(
                `Search failed: ${error instanceof Error ? error.message : 'Unknown error'}`
            );
        }
    }

    async read<T>(model: string, ids: number[], fields?: string[]): Promise<T[]> {
        const url = `${this.config.url}${model}/read`;

        try {
            const response = await fetch(url, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    "Authorization": `Bearer ${this.config.apikey}`,
                },
                body: JSON.stringify({ ids: ids, fields: fields }),
            });

            if (!response.ok) throw new Error(`HTTP error: ${response.status}`);

            const data = await response.json() as T[];
            return data;
        } catch (error) {
            throw new OdooRPCError(
                `Read failed: ${error instanceof Error ? error.message : 'Unknown error'}`
            );
        }
    }

    async search_read<T>(
        model: string,
        domain?: string,
        fields?: string[],
        offset?: number,
        limit?: number
    ): Promise<T[]> {
        const url = `${this.config.url}${model}/search_read`;
        try {
            const response = await fetch(url, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    "Authorization": `Bearer ${this.config.apikey}`,
                },
                body: JSON.stringify({
                    domain: parseDomainString(domain ?? ""),
                    fields: fields,
                    offset: offset ?? 0,
                    limit: limit ?? 0,
                }),
            });

            if (!response.ok) throw new Error(`HTTP error: ${response.status}`);

            const data = await response.json() as T[];
            return data;
        } catch (error) {
            throw new OdooRPCError(
                `Search_read failed: ${error instanceof Error ? error.message : 'Unknown error'}`
            );
        }
    }

    async write(model: string, id: number | number[], data: any): Promise<any> {
        const url = `${this.config.url}${model}/write`;

        try {
            const response = await fetch(url, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    "Authorization": `Bearer ${this.config.apikey}`,
                },
                body: JSON.stringify({ id: id, data: data }),
            });

            if (!response.ok) throw new Error(`HTTP error: ${response.status}`);

            const result = await response.json();
            return result;
        } catch (error) {
            throw new OdooRPCError(
                `Write failed: ${error instanceof Error ? error.message : 'Unknown error'}`
            );
        }
    }

    async unlink(model: string, ids: number[]): Promise<any> {
        const url = `${this.config.url}${model}/unlink`;

        try {
            const response = await fetch(url, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    "Authorization": `Bearer ${this.config.apikey}`,
                },
                body: JSON.stringify({ ids: ids }),
            });

            if (!response.ok) throw new Error(`HTTP error: ${response.status}`);

            const result = await response.json();
            return result;
        } catch (error) {
            throw new OdooRPCError(
                `Unlink failed: ${error instanceof Error ? error.message : 'Unknown error'}`
            );
        }
    }

    async execute(model: string, method: string, args: any): Promise<any> {
        const url = `${this.config.url}${model}/execute`;

        try {
            const response = await fetch(url, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    "Authorization": `Bearer ${this.config.apikey}`,
                },
                body: JSON.stringify({ method: method, args: args }),
            });

            if (!response.ok) throw new Error(`HTTP error: ${response.status}`);

            const result = await response.json();
            return result;
        } catch (error) {
            throw new OdooRPCError(
                `Execute failed: ${error instanceof Error ? error.message : 'Unknown error'}`
            );
        }
    }

    async execute_kw(model: string, method: string, args: any[] = [], kwargs: any = {}): Promise<any> {
        const url = `${this.config.url}${model}/execute_kw`;

        try {
            const response = await fetch(url, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    "Authorization": `Bearer ${this.config.apikey}`,
                },
                body: JSON.stringify({ method: method, args: args, kwargs: kwargs }),
            });

            if (!response.ok) throw new Error(`HTTP error: ${response.status}`);

            const result = await response.json();
            return result;
        } catch (error) {
            throw new OdooRPCError(
                `Execute_kw failed: ${error instanceof Error ? error.message : 'Unknown error'}`
            );
        }
    }
}

export interface JRPCConfig {
    hostname: string;
    port: number;
    schema: string;
    database: string;
    username: string;
    password: string;
    apikey: string;
    url?: string;
    uid?: number;
}

export class JRPCClient {
    constructor(private strategy: JRPC) { }

    async login(): Promise<number> {
        return this.strategy.login();
    }

    async create(model: string, data: any): Promise<any> {
        return this.strategy.create(model, data);
    }

    async load(model: string, header: any, values: any): Promise<any> {
        return this.strategy.load(model, header, values);
    }

    async count(model: string, domain?: string): Promise<number> {
        return this.strategy.count(model, domain);
    }

    async fields_get(model: string, fields?: string[], field_attributes?: string[]): Promise<any> {
        return this.strategy.fields_get(model, fields, field_attributes);
    }

    async get_id(model: string, domain?: string): Promise<number> {
        return this.strategy.get_id(model, domain);
    }

    async search(model: string, domain?: string): Promise<number[]> {
        return this.strategy.search(model, domain);
    }

    async read<T>(model: string, ids: number[], fields?: string[]): Promise<T[]> {
        return this.strategy.read(model, ids, fields);
    }

    async search_read<T>(
        model: string,
        domain?: string,
        fields?: string[],
        offset?: number,
        limit?: number
    ): Promise<T[]> {
        return this.strategy.search_read(model, domain, fields, offset, limit);
    }

    async write(model: string, id: number | number[], data: any): Promise<any> {
        return this.strategy.write(model, id, data);
    }

    async unlink(model: string, ids: number[]): Promise<any> {
        return this.strategy.unlink(model, ids);
    }

    async execute(model: string, method: string, args: any): Promise<any> {
        return this.strategy.execute(model, method, args);
    }

    async execute_kw(model: string, method: string, args: any[] = [], kwargs: any = {}): Promise<any> {
        return this.strategy.execute_kw(model, method, args, kwargs);
    }

    // Optional: allow runtime strategy swap
    setStrategy(strategy: JRPC): void {
        this.strategy = strategy;
    }
}

export type Protocol = "jsonrpc" | "odoojson";

export class JRPCStrategyFactory {
    static create(protocol: Protocol, config: JRPCConfig): JRPC {
        switch (protocol) {
            case "jsonrpc":
                return new JSONRPC(config);

            case "odoojson":
                return new ODOOJSON(config);

            default:
                throw new Error(`Unsupported protocol: ${protocol}`);
        }
    }
}


// Protocol-specific extras
function parseDomainString(input: string): (string | string[])[] {
    // Remove the outer [ and ]
    let cleaned = input.trim();
    if (cleaned === "") return [];
    if (cleaned.startsWith('[')) cleaned = cleaned.slice(1);
    if (cleaned.endsWith(']')) cleaned = cleaned.slice(0, -1);

    const result: (string | string[])[] = [];

    // Check if there's an operator at the start (like '|' or '&')
    const operatorMatch = cleaned.match(/^'([|&])'\s*,\s*/);
    if (operatorMatch) {
        result.push(operatorMatch[1]!);
        // Remove the operator from the string
        cleaned = cleaned.slice(operatorMatch[0].length);
    }

    // Split by "),(" to separate each tuple
    const tuples = cleaned.split('),(');

    tuples.forEach(tuple => {
        // Remove any remaining parentheses
        const withoutParens = tuple.replace(/[()]/g, '');

        // Split by comma and clean up quotes/spaces
        const parsed = withoutParens
            .split(',')
            .map(item => item.trim().replace(/^'|'$/g, ''));

        result.push(parsed);
    });

    return result;
}