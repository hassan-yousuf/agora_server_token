package main

import (
    "net/http"
    "os"
    "strconv"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    rtctokenbuilder "github.com/AgoraIO-Community/go-tokenbuilder/rtctokenbuilder"
    rtmtokenbuilder "github.com/AgoraIO-Community/go-tokenbuilder/rtmtokenbuilder"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        panic("Error loading .env file")
    }

    appID := os.Getenv("APP_ID")
    appCertificate := os.Getenv("APP_CERTIFICATE")
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    router := gin.Default()

    router.GET("/rtcToken", func(c *gin.Context) {
        channelName := c.Query("channelId")
        uidStr := c.Query("uid")
        uid, _ := strconv.Atoi(uidStr)
        role := rtctokenbuilder.RolePublisher
        expireTimeInSeconds := uint32(3600)
        currentTimestamp := uint32(time.Now().Unix())
        privilegeExpiredTs := currentTimestamp + expireTimeInSeconds

        token, err := rtctokenbuilder.BuildTokenWithUID(appID, appCertificate, channelName, uint32(uid), rtctokenbuilder.Role(role), privilegeExpiredTs)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, gin.H{"rtcToken": token})
    })

    router.GET("/rtmToken", func(c *gin.Context) {
        userAccount := c.Query("userAccount")
        expireTimeInSeconds := uint32(3600)
        currentTimestamp := uint32(time.Now().Unix())
        privilegeExpiredTs := currentTimestamp + expireTimeInSeconds

        token, err := rtmtokenbuilder.BuildToken(appID, appCertificate, userAccount, rtmtokenbuilder.RoleRtmUser, privilegeExpiredTs)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, gin.H{"rtmToken": token})
    })

    router.Run(":" + port)
}
