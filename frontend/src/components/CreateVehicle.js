import React, { useState, useEffect } from 'react'
import {createVehicle} from "../api";

export default function CreateVehicle({history}) {
    const [error, setError] = useState("")
    const [year, setYear] = useState(2000)
    const [make, setMake] = useState("Ford")
    const [model, setModel] = useState("Focus")

    const submit = async () => {
        try {
            const vehicle = await createVehicle(parseInt(year), make, model)
            history.push('/vehicles/' + vehicle.vehicle_id)
        } catch (e) {
            setError(e.toString())
        }

        return false
    }

    return (
        <form>
            <h3>Create A Vehicle</h3>
            {error ? <p className="error">{error}</p> : null}

            <div className="control">
                <label htmlFor="year">Year</label>
                <input type="text" onChange={e => setYear(e.target.value)} />
            </div>

            <div className="control">
                <label htmlFor="make">Make</label>
                <input type="text" onChange={e => setMake(e.target.value)} />
            </div>

            <div className="control">
                <label htmlFor="model">Model</label>
                <input type="text" onChange={e => setModel(e.target.value)} />
            </div>

            <button type="button" onClick={submit}>Create Vehicle</button>
        </form>
    )
}
