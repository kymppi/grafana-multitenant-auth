const http = require('http');

const tenants = [
  {
    id: '1',
  },
  {
    id: '2',
  },
];

const server = http.createServer((req, res) => {
  console.log(req.url, req.method, req.headers);

  if (req.url === '/auth/') {
    const tenantId = req.headers['x-original-uri'].split('/')[1];

    const tenant = tenants.find((t) => t.id === tenantId);

    if (tenant) {
      res.writeHead(200, {
        'Content-Type': 'text/plain',
        'X-Scope-OrgId': tenant.id,
      });
      res.end('OK\n');
    } else {
      res.writeHead(401, {
        'Content-Type': 'text/plain',
      });
      res.end('Not Found\n');
    }
  } else if (req.url === '/api/v1/push') {
    // get the tenant id from the request
    const tenantId = req.headers['x-scope-orgid'];

    // get the tenant from the list of tenants
    const tenant = tenants.find((t) => t.id === tenantId);

    // if the tenant exists
    if (tenant) {
      // send the response
      res.writeHead(200, {
        'Content-Type': 'text/plain',
      });
      res.end('OK\n');
    } else {
      // if the tenant does not exist
      res.writeHead(401, {
        'Content-Type': 'text/plain',
      });
      res.end('Not Found\n');
    }
  } else {
    res.writeHead(404, {
      'Content-Type': 'text/plain',
    });

    res.end('Not Found\n');
  }
});

const PORT = 3001;
server.listen(PORT, () => {
  console.log(`Server listening on port ${PORT}`);
});
