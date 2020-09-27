const URL = 'http://localhost:8000'


export async function isAuthenticated() {
    const response = await fetch(
        URL + '/v1/session'
    )
}

class FailedApiRequest {
    constructor(response, body) {
        this.response = response
        this.body = body
    }

    toString() {
        return JSON.stringify(this.body)
    }
}

async function request(method, path, body) {
    const options = {
        'method': method,
        'headers': {},
        'credentials': 'include',
    }

    const headers = options.headers
    headers['content-type'] = headers['content-type'] || 'application/json'
    headers['accept'] = headers['accept'] || 'application/json'

    if (body) {
        headers['content-type'] = 'application/json'
        options['body'] = JSON.stringify(body)
    }

    const response = await fetch(URL + path, options)

    if (response.status === 204) {
        return null
    }

    const responseBody = await response.json()
    if (response.status >= 400) {
        throw new FailedApiRequest(response, responseBody)
    }

    return responseBody
}

export async function createUser(emailAddress, password) {
    const response = await request(
        'POST', '/v1/users/',
        {
            email_address: emailAddress,
            password,
        })

    return response
}

export async function getSession() {
    const response = await request('GET', '/v1/session')
    return response
}

export async function createSession(emailAddress, password) {
    await request(
        'POST', '/v1/session', {email_address: emailAddress, password},
    )
}

export async function  deleteSession() {
    await request(
        'DELETE', '/v1/session',
    )
}

export async function getVehicles() {
    return await request('GET', '/v1/vehicles/')
}

export async function createVehicle(year, make, model) {
    return await request('POST', '/v1/vehicles/', {year, make, model})
}

export async function deleteVehicle(vehicleId) {
    return await request('DELETE', `/v1/vehicles/${vehicleId}`)
}
