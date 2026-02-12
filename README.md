# odoojrpc

A TypeScript/JavaScript library for interacting with Odoo via JSON-RPC and JSON2 protocols.

## Installation

```bash
npm install @ppreeper/odoojrpc
```

```bash
npm install git+https://github.com/ppreeper/odoojrpc.git
```

## Features

- 🚀 Support for both JSONRPC and ODOOJSON (JSON2) protocols
- 💪 Full TypeScript support with type definitions
- 🎯 Strategy pattern for easy protocol switching
- ✅ Comprehensive error handling
- 🔒 Type-safe operations

## Quick Start

### JSONRPC Protocol

```typescript
import { JRPCClient, JRPCStrategyFactory } from '@ppreeper/odoojrpc';

const config = {
    hostname: 'localhost',
    port: 8069,
    schema: 'http',
    database: 'mydb',
    username: 'admin',
    password: 'admin',
    apikey: '' // Not used for JSONRPC
};

const client = new JRPCClient(
    JRPCStrategyFactory.create('jsonrpc', config)
);

// Login
const uid = await client.login();
console.log('Logged in as user:', uid);

// Search and read records
const partnerIds = await client.search('res.partner', [['is_company', '=', true]]);
const partners = await client.read('res.partner', partnerIds, ['name', 'email']);

// Create a record
const newPartnerId = await client.create('res.partner', {
    name: 'New Company',
    email: 'info@newcompany.com'
});

// Update a record
await client.write('res.partner', newPartnerId, {
    phone: '+1234567890'
});

// Delete a record
await client.unlink('res.partner', [newPartnerId]);
```

### ODOOJSON (JSON2) Protocol

```typescript
import { JRPCClient, JRPCStrategyFactory } from '@ppreeper/odoojrpc';

const config = {
    hostname: 'localhost',
    port: 8069,
    schema: 'http',
    database: 'mydb',
    username: '', // Not used for ODOOJSON
    password: '', // Not used for ODOOJSON
    apikey: 'your-api-key-here'
};

const client = new JRPCClient(
    JRPCStrategyFactory.create('odoojson', config)
);

// No login required for ODOOJSON - authentication via API key
const partners = await client.search_read('res.partner',
    [['is_company', '=', true]],
    ['name', 'email']
);
```

## API Reference

### Client Methods

#### `login(): Promise<number>`
Authenticate with the Odoo server (JSONRPC only). Returns the user ID.

#### `search(model: string, domain?: any[]): Promise<number[]>`
Search for records matching the domain. Returns array of record IDs.

#### `read<T>(model: string, ids: number[], fields?: string[]): Promise<T[]>`
Read specific records by ID.

#### `search_read<T>(model: string, domain?: any[], fields?: string[], offset?: number, limit?: number): Promise<T[]>`
Combined search and read operation.

#### `create(model: string, data: any): Promise<any>`
Create a new record.

#### `write(model: string, id: number | number[], data: any): Promise<any>`
Update existing record(s).

#### `unlink(model: string, ids: number[]): Promise<any>`
Delete record(s).

#### `count(model: string, domain?: any[]): Promise<number>`
Count records matching the domain.

#### `fields_get(model: string, fields?: string[], attributes?: string[]): Promise<any>`
Get field definitions for a model.

#### `execute(model: string, method: string, args: any): Promise<any>`
Execute a model method.

#### `execute_kw(model: string, method: string, args?: any[], kwargs?: any): Promise<any>`
Execute a model method with keyword arguments.

## TypeScript Support

The library is written in TypeScript and provides full type definitions:

```typescript
import type { JRPCConfig, JRPC, Protocol } from '@ppreeper/odoojrpc';

// Define your own model interfaces
interface Partner {
    id: number;
    name: string;
    email: string;
}

const partners = await client.read<Partner>('res.partner', [1, 2, 3], ['name', 'email']);
// partners is typed as Partner[]
```

## Error Handling

The library throws `OdooRPCError` for all RPC-related errors:

```typescript
import { OdooRPCError } from '@ppreeper/odoojrpc';

try {
    await client.login();
} catch (error) {
    if (error instanceof OdooRPCError) {
        console.error('Odoo error:', error.message);
        console.error('Error code:', error.code);
        console.error('Error data:', error.data);
    }
}
```

## Switching Protocols at Runtime

```typescript
const client = new JRPCClient(
    JRPCStrategyFactory.create('jsonrpc', config)
);

// Later, switch to ODOOJSON
client.setStrategy(JRPCStrategyFactory.create('odoojson', config));
```

## Domain Syntax

Odoo uses a specific domain syntax for filtering:

```typescript
// Simple equality
[['name', '=', 'John']]

// Multiple conditions (AND)
[['is_company', '=', true], ['country_id', '=', 1]]

// OR conditions
['|', ['name', '=', 'John'], ['name', '=', 'Jane']]

// Comparison operators
[['id', '>', 100]]
[['create_date', '>=', '2024-01-01']]

// LIKE operator
[['name', 'ilike', 'john']]  // case-insensitive

// IN operator
[['id', 'in', [1, 2, 3]]]
```

## License

BSD-2-Clause

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Author

Peter Preeper