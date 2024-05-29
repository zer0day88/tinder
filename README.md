### How to test the project
run `make test`

> for the test I use integration test by using https://testcontainers.com/

### How to run this service

1. install docker & docker-compose
2. run `make docker-up`
3. if you're using different container tool such as podman <br>
   you can replace docker-compose with `docker compose` or `podman compose`
4. the command will run postgres, redis and the service on port 3000

### This service contains 3 APIs:
1. POST `/v1/signup` (to register user) <br>
   payload <br>
   ```json
   {
      "email": "johndoe@gmail.com",
      "password": "Abcdefghij$1"
   }
   ```
   success response
   ```json
   {
      "code": 200,
      "message": "success"
   }
   ```
2. POST `/v1/login` (to login) <br>
   payload
   ```json
   {
      "email": "johndoe@gmail.com",
      "password": "Abcdefghij$1"
   }
   ```
   success response
   ```json
   {
    "code": 200,
    "message": "Success",
    "user": {
        "AccessToken": {
            "Token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTY5NTA2MTMsImlhdCI6MTcxNjk0ODgxMywibmJmIjoxNzE2OTQ4ODEzLCJzdWIiOiI4ODk4MDlkZi1mMDNlLTQzOTMtODdmZS03NzE2YWJhYjcxZDUiLCJ0b2tlbl91dWlkIjoiMjIzNjRjNjUtMmJhNS00NDFlLTg4YzMtZmU3MmM1OTYxYWE5In0.LMAdfYLmjjfmyVFAneTScgMOBSyxqF4eRYWv3jzhuzG6Ge0A-nM2fTzz4mdZb0m11_fc-D4kX-9bBvyx_5HOx6naCu-mEjX7rOOeDyAc4Oriwr-Il5Jegkdrp-uY0-RXSvlLexb-3Mdkgr1aDCkcITZETN9M3tMpDAPv4n75zh7qNKiuDRp8Bw8Fuwf9MgimGI3JNKZLwqOtIUTmNA2ZRSmXr7hDDkT-SYhmW5mB4GJSkYTSno_T0wdQOekCF-It-zGczsKE5BlM91iIzB_fSgkbpH5HVlmNGW4ku-C2SZiDOvVo0jp-IcxYBs_MOEzk3OJ-bsZZjg-uUFj0gyTtnndAMfemlVKtzY4WzK1wg0GHSD6Td5DDmfjpZVhN0Ojs9HLPFoCmFxug_-FLwy49KjMCZAil-W_UQhQzhchlThKvg_S4WtU0DWbKiOgrvB8pMu8kf77PWBcH-1eq_IQ5B1PfBeLe3NnLqllvu49wPl6yvnRIncfNhpTiKwYHpklNCbPwbA4EriCSw4WxAgsDk4CpGUhfBHT960OVV0WlXxDi4ZsSJbzVO0hBpVZtYCzbD_vIz4ym_3TMXS2D6blP9dmKMYQ6sO6TV0f6FARkbJSMqPoW83lNCy-mAAAqcylxJhRoFaPh2I7L2u1XcNwQDr6bw-FKzagbbArRXhRIOaM",
            "TokenUuid": "22364c65-2ba5-441e-88c3-fe72c5961aa9",
            "UserID": "889809df-f03e-4393-87fe-7716abab71d5",
            "ExpiresIn": 1716950613
        },
        "RefreshToken": {
            "Token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTY5NTI0MTMsImlhdCI6MTcxNjk0ODgxMywibmJmIjoxNzE2OTQ4ODEzLCJzdWIiOiI4ODk4MDlkZi1mMDNlLTQzOTMtODdmZS03NzE2YWJhYjcxZDUiLCJ0b2tlbl91dWlkIjoiZTA3MWZiM2UtNTI2YS00ZjJlLTg2MTItMDQ0YjkxOTdiZjE3In0.JMiGk8ma1JmCPl1uVNxnmoESv6Hsxs8j8_i2vbwSWcIYTdwxJqdCIQ7UDiWj9NJXAY7T1yuXcbwdjrOZalyHckp_dbyEK9u9uw6gwrYdp6UBATNI5ylxXBOuUWiRIVbIrtAHTsjkZ-NvEtLF7wADX0_Ah0Ippw8iVM9rEeaJGNC7RyUSyl0KVmfW-l6B4Apl4yatk-JouGqDW7B_sxU9JIyNnHmQgN-5GrAmMLJQCrqp49x6oc7_2Z3-xXQoawpmcRE6DjyVRJjIgTzVyt2U03Le62X1EfrVhIzPLhNv2eniJWMc64NWY6rl0-zEaCcQY6lShn-XN7D3q4kdKPH4loyCuISxN9DINAc8U1sOnVc0-RQzS1QX20f-F--ZWf9nhUzwhOnYeN2eRNgeawipu3lT2xWrAgdbuSguLwh2zPo6oZYazQ33N3duPuin21L6PkQs2_h_zxtWn-F3gD5k0JXBmd5PepjiGB-NIxjFplmB1Zkqt7UfItj_lgZMtBWVN3e6GUMhfwe3JRYpip1BBFOyBQ8LKUVht9sXxBGtOxubfp_NB-mQbRodDHVRxkGj0uL8E6f2o-B7IPf5VkKFsbOLacWReOWWGTahbuHiefExc67uexWjN4vCvToAcPe6hhHuyC-Ndkp_OFY_CzFJYtSDntlKBL3t43eAtjn52fQ",
            "TokenUuid": "e071fb3e-526a-4f2e-8612-044b9197bf17",
            "UserID": "889809df-f03e-4393-87fe-7716abab71d5",
            "ExpiresIn": 1716952413
         }
      }
   }
   ```
3. GET `/v1/cek` (to check authentication) <br>
   > use access token from login response
   
   header
   `Authorization: Bearer <access token>`
   
   success response
   ```json
   {
      "code": 200,
      "message": "success"
   }
   ```
   failed response
   ```json
   {
      "code": 401,
      "message": "Request unauthorized. Missing or invalid access token"
   }
   ```

### The service structure:
1. `cmd` contains the bootstrap file for running the service <br>
   from load the config,db..etc and graceful shutdown the services
2. `config` contains the configuration file for the service
3. `internal` this is isolated layer that contains the route, handle, services, repository, middleware, etc
4. `migrations` contain migration file and the migration config as well
5. `pkg` package for general purpose
6. `seed` file to run seeder
7. `infra` contains the infrastructure package
8. `helper` contain utility to help the main service