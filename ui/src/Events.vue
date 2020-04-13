<template>
    <div class="events-wrapper">
        <div v-for="e in events">
            <div class="event"
                 style="display: flex;"
                :style="e.style">
                <div class="event-field">{{ e.service }}</div>
                <!-- TODO: Include VueMQ
                <!--<mq-layout mq="md+">-->
                    <div class="event-field">{{ e.title }}</div>
                <!--</mq-layout>-->
                <button @click="deleteEvent(e.id)">&times;</button>
            </div>
            <div v-show="false"> more info here </div>
        </div>
    </div>
</template>

<script>
    import api from "./requests.js";
    export default {
        name: "Events",
        data() {
            return {
                loadedEvents: [],
                statusStyling: {
                    "ok": { color: "#000", backgroundColor: "#c3e6cb" },
                    "warning": { color: "#000", backgroundColor: "#ffeeba" },
                    "error": { color: "#000", backgroundColor: "#f5c6cb" },
                    "info": { color: "#000", backgroundColor: "#fff" }
                }
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
            deleteEvent(id) {
                api.deleteEvent(id, 
                    success => {
                        this.loadedEvents = this.loadedEvents.filter(
                            e => e.id !== id
                        );
                    },
                    error => {
                        console.error(error)
                        this.$emit("error");
                    }
                ); 
            }
        },
        computed: {
            events() {
                return this.loadedEvents.map(e => {
                    return {
                        ...e,
                        style: this.statusStyle(e.status)
                    };
                });
            }
        },
        mounted() {
            api.getEvents(
                successData => {
                    this.loadedEvents = successData;
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
    overflow-x: scroll;
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
