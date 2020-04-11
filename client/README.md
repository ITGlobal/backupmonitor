# BackupMonitor: Frontend

## Local development

1. Start backend server at `http://localhost:8000`:

   ```shell
   make run
   
   # or

   go build
   ./backupmonitor
   ```

2. Start frontend development server:

   ```shell
   cd client
   npm install
   npm run start
   ```

3. Navigate to `http://localhost:4200/`.

## Production build

Type the following commands:

```shell
cd client
npm install
npm run build
```

Output is packaged to `www` directory located in repository root.
