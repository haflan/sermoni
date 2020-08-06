<template>
    <div class="events-wrapper">
        <div v-for="e in events" :key="e.id">
            <div class="event"
                 style="display:flex;cursor:pointer;"
                :style="e.style"
                @click="eventClicked(e, $event)">
                <div class="event-field">{{ e.service }}</div>
                <!-- TODO: Include VueMQ
                <!--<mq-layout mq="md+">-->
                    <div class="event-field">{{ simplifyDate(e.timestamp) }}</div>
                <!--</mq-layout>-->
                <button v-show="e.id > 0" class="delete-button">&times;</button>
            </div>
            <div v-show="e.open" style="padding: 5px; border: 1px solid rgba(0,0,0,.125); border-top: 0">
                <code style="white-space:pre;">{{e.details}}</code>
            </div>
        </div>
    </div>
</template>

<script>
    import api from "./requests.js";
    export default {
        name: "Events",
        data() {
            return {
                events: [],
                statusStyling: {
                    "ok": { color: "black", backgroundColor: "#c3e6cb" },
                    "warning": { color: "black", backgroundColor: "#ffeeba" },
                    "error": { color: "black", backgroundColor: "#f5c6cb" },
                    "info": { color: "black", backgroundColor: "#fff" },
                    "late":{ color: "black", backgroundColor: "#d6d8d9" },
                },
                eventSorter: (e1, e2) => e2.timestamp - e1.timestamp
            };
        },
        methods: {
            statusStyle(status) {
                const style = this.statusStyling[status];
                if (style) {
                    return style;
                } else {
                    return { color: "#000", backgroundColor: "#fff" };
                }
            },
            eventClicked(e, clickEvent) {
                if (clickEvent.target.className === "delete-button") {
                    // Element clicked is the delete button
                    this.deleteEvent(e.id);
                } else {
                    // Not the delete button - toggle details
                    e.open = !e.open;
                }
            },
            deleteEvent(id) {
                api.deleteEvent(id, 
                    success => {
                        this.events = this.events.filter(
                            e => e.id !== id
                        );
                    },
                    error => {
                        console.error(error);
                        this.$emit("error");
                    }
                ); 
            },
            simplifyDate(eventUTC) {
                let diffUTC = Date.now() - eventUTC;
                let inFuture = false;
                if (diffUTC < 0) {
                    inFuture = true;
                    diffUTC = -diffUTC;
                }
                let unitStrings = ["year", "week", "day", "hour", "minute"];
                let units = [31536000000, 604800000, 86400000, 3600000, 60000];
                let numOfUnits = 0;
                for (let i = 0; i < units.length; i++) {
                    numOfUnits = Math.floor(diffUTC / units[i]);
                    if (numOfUnits === 0) {
                        continue;
                    }
                    if (inFuture) {
                        return "in " + numOfUnits + " " + unitStrings[i] + (numOfUnits > 1 ? "s" : "");
                    }
                    return numOfUnits + " " + unitStrings[i] + (numOfUnits > 1 ? "s" : "") + " ago";
                }
                return inFuture ? "now" : "just now";
            }
        },
        mounted() {
            api.getEvents(
                successData => {
                    let lateId = 0;
                    this.events = successData.map(e => {
                        return {
                            ...e,
                            id: e.id ? e.id : lateId--,
                            style: this.statusStyle(e.status),
                            open: false
                        }
                    }).sort(this.eventSorter);
                },
                error => {
                    console.error(error);
                    this.$emit("error");
                }
            );
        }
    }

</script>

<style scoped>
.event {
    height: 3em;
    /*
    margin-bottom: 0.5em;
    -webkit-box-shadow: 0 0 10px rgba(0,0,0,.1);
    box-shadow: 0 0 10px rgba(0,0,0,.1);
    */
    padding: .75rem;
    border: 1px solid rgba(0,0,0,.125);
    box-sizing: border-box;
    white-space: nowrap;
}
.event-field {
    flex: 1;
}
.events-wrapper {
    margin: 0;
    padding: 0;
}
button {
    cursor: pointer;
    background-color: inherit;
    border: none;
    text-align: center;
    font-size: 16px;
}
</style>
