export async function getVehicles() {
    const response = await fetch(
        'http://localhost:8000/v1/vehicles/',
        {
            method: 'GET',
        },
    )

    const body = await response.json()
    console.log("3", body)

    return []
}