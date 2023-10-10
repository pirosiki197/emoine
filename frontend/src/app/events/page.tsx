'use client'

import { APIService } from "@/api/v1/api_connect"
import { Event } from "@/api/v1/api_pb"
import { useClient } from "@/use-client"
import { ConnectError } from "@connectrpc/connect"
import { useState } from "react"

function Events() {
    const [events, setEvents] = useState<Event[]>([])
    const res = useClient(APIService).getEvents({request: {}})
    res.then((res) => {
        setEvents(res.events)
    }).catch((err) => {
        const connectError = ConnectError.from(err)
        console.log(connectError.code, connectError.message)
    })

    return (
        <div>
            <ul>
                {events.map((event) => {
                    return (
                        <li key={event.id}>
                            {event.title}
                        </li>
                    )
                })}
            </ul>
        </div>
    )
}

export default function Page() {
    return (
        <div>
            <h1>Events</h1>
            <Events />
        </div>
    )
}