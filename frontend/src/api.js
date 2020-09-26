const URL = 'http://localhost:8000'


export async function isAuthenticated() {
    const response = await fetch(
        URL + '/v1/session'
    )
}

export async function createUser(emailAddress, password) {
    const response = await fetch(
        '/v1/users/',
        {
            method: 'POST',
            headers: {
                'content-type': 'application/json',
            },
            body: JSON.stringify({
                email_address: emailAddress,
                password,
            }),
        },
    )

    const body = await response.json()
    if (response.status >= 400) {
        throw body
    }

    return body
}

export async function getVehicles() {
    const response = await fetch(
        URL + '/v1/vehicles/',
        {method: 'GET'},
    )

    const body = await response.json()
    return body
}