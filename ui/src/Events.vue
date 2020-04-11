<template>
    <div class="events-wrapper">
        <div v-for="e in events">
            <div class="event"
                 style="display: flex;"
                :style="e.style">
                <div class="event-field">{{ e.service }}</div>
                <!-- TODO: Include VueMQ
                <mq-layout mq="md+">
                    <div class="event-field">{{ e.title }}</div>
                </mq-layout>
                -->
            </div>
            <div v-show="false"> more info here </div>
        </div>
    </div>
</template>

<script>
    export default {
        name: "Events",
        data() {
            return {
                testdata: [{
                    service: "backup-files @ haflan.dev",
                    timestamp: 1586459792925,
                    title: "Server reported success",
                    status: "ok"
                },{
                    service: "backup-gitlab @ haflan.dev",
                    timestamp: 1586459793155,
                    title: "Server reported success",
                    status: "ok"
                },{
                    service: "backup-qumps @ haflan.dev",
                    timestamp: 1586459793285,
                    status: "warning",
                    title: "Memory almost full",
                    message: "df -h reports that less than 1GB is available"
                },{
                    service: "backup-offsite @ work-computer",
                    timestamp: 1586459794385,
                    status: "error",
                    title: "Expectation not met" // These are read from a default title thingy
                },{
                    service: "ssh @ haflan.dev",
                    timestamp: 1586459794385,
                    status: "info",
                    title: "SSH server login",
                    message: "User vetle logged in from IP 192.168.10.105"
                }],
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
                    return { color: "#fff", backgroundColor: "#000" };
                }
            }
        },
        computed: {
            events() {
                return this.testdata.map(e => {
                    return {
                        ...e,
                        style: this.statusStyle(e.status)
                    };
                });
            }
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
</style>
