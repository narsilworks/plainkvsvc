
printf "Building binary..."
go build ./
printf "done.\r\n"

for n in {1..3}
do
    # Stop existing services
    printf -v svc "plainkvsvc-%s.service" $n
    if [[ $(systemctl list-units --all -t service --full --no-legend "$svc" | cut -f1 -d' ') == $svc ]]; then
        printf "Stopping service $svc..."
        systemctl stop "$svc"
        printf "done.\r\n"
    fi

    printf "Generating new service..."

    # create root directory
    rd="/usr/local/services/plainkvsvc/api"
    mkdir -p "$rd"
    chown -R plainkvsvc:plainkvsvc "$rd"
    chmod +755 "$rd"

    # Copy existing services
    printf -v dest "$rd/%s" $n

    # Create directory if it does not exist
    mkdir -p "$dest"
    chown -R plainkvsvc:plainkvsvc "$dest"
    chmod +755 "$dest"

    # Copy api to destination
    cp api "$dest"

    # Copy config if it does not exist
    #cp config.json "$dest"

    # Copy localdb
    cp local.db "$dest"
    chmod +775 "$dest/local.db"

    a="plainkvsvc-$n.service"

    # Copy service config
    cp plainkvsvc-0.service "$a"

    # Replace service variables
    sed -i "s/{0}/$n/g" "$a" # change api number
    sed -i "s/{1}/$n/g" "$a" # change api port

    aa="/etc/systemd/system/$a"

    # Move to service
    if test -f "$aa"; then
        rm -rf "$aa"
    fi

    mv "$a" "$aa"
    rm -rf "$a"

    printf "done.\r\n"

    # Reload systemctl
    systemctl daemon-reload

    # Start service
    printf "Starting service $a..."
    systemctl enable "$a"
    systemctl start "$a"
    printf "done.\r\n"
done