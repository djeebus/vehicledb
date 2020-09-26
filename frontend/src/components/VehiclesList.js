import React, { useState, useEffect } from 'react'
import { getVehicles } from "../api";

export default function() {
    const [vehicles, setVehicles] = useState([])
    console.log("1", vehicles)

    useEffect(() => {
        async function loadVehicles() {
            console.log("3")
            const cars = await getVehicles()
            console.log("2", cars)
            setVehicles(cars)
        }
        loadVehicles()
    }, [])

    return <div>{ vehicles.map(v => v) }</div>
}
