import React, { useState, useEffect } from 'react'
import ReactDOM from 'react-dom'
import { getVehicles } from "./api";

const App = () => {
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

const app = document.getElementById('app')
ReactDOM.render(<App />, app)