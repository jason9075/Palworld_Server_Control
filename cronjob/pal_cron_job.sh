#!/bin/bash

# 設定檔案來儲存連續檢查的次數
COUNT_FILE="/tmp/player_count.txt"
SERVER_COMMAND="/home/jason/.local/bin/stop_pal.sh" 
SERVER_ADDRESS="xx.xx.xx.xx:8211"
DISCORD_WEBHOOK="https://discord.com/api/webhooks/xxx/xxxxx"
DOWN_MSG="帕爾勞工們已下班..."

# 執行gamedig並獲取numplayers
NUM_PLAYERS=$(gamedig --type palworld --host "$SERVER_ADDRESS" | jq '.numplayers')

# write to log file
echo "$(date "%D %T") - $NUM_PLAYERS" >> /tmp/pal_cron_job.log

# check null
if [ "$NUM_PLAYERS" == "null" ]; then
    exit 0
fi

# 檢查檔案是否存在，如果不存在則創建
if [ ! -f "$COUNT_FILE" ]; then
    echo 0 > "$COUNT_FILE"
fi

# 獲取當前連續為0的次數
COUNT=$(cat "$COUNT_FILE")

# 檢查玩家數量
if [ "$NUM_PLAYERS" -eq 0 ]; then
    # 增加計數
    COUNT=$((COUNT + 1))
    echo $COUNT > "$COUNT_FILE"

    # 檢查是否連續三次為0
    if [ "$COUNT" -ge 3 ]; then
        curl -H "Content-Type: application/json" -X POST -d "{\"content\": \"$DOWN_MSG\"}" "$DISCORD_WEBHOOK"
        # 執行關閉伺服器命令
        $SERVER_COMMAND
        # 重置計數
        echo 0 > "$COUNT_FILE"
        systemctl suspend
    fi
else
    # 如果玩家數量不為0，重置計數
    echo 0 > "$COUNT_FILE"
fi
