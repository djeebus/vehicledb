import React, { useState, useEffect } from 'react'
import { getVehicles, deleteVehicle } from "../api";
import CreateVehicle from "./CreateVehicle";

export default function({history}) {
    const [vehicles, setVehicles] = useState([])

    useEffect(() => {
        async function loadVehicles() {
            const cars = await getVehicles()
            setVehicles(cars)
        }
        loadVehicles()
    }, [])

    const onDelete = vehicle => {
        deleteVehicle(vehicle.vehicle_id)
    }

    function renderVehicle(vehicle) {
        return (
            <li key={vehicle.vehicle_id}>
                <a href={`/vehicles/${vehicle.vehicle_id}`}>
                    {vehicle.year} {vehicle.make} {vehicle.model}
                </a>

                (<a href="" onClick={onDelete(vehicle)}>Delete me</a>)
            </li>
        )
    }

    return (
        <div>
            <ul>
                { vehicles.map(renderVehicle) }
            </ul>
            <CreateVehicle history={history} />
        </div>
    )
}
