import { resolve } from 'path';

import { generateApi } from 'swagger-typescript-api';

generateApi({
    name: 'api.ts',
    output: resolve(process.cwd(), './src/api'),
    url: 'https://localhost:8000/api/docs/doc.json',
    httpClientType: 'axios',
});