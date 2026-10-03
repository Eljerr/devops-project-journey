const request = require('supertest')
const app = require('../app')

describe('Sanity Check Api', () => {
  it('seharusnya mengembalikan HTTP 200 saat mengakses root endpoint', async () => {
      const res = await request(app).get('/')
      expect(res.statusCode).toEqual(200)
  })
})
