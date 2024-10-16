namespace ATCLib.Test;

using Xunit;

public class MessageTests
{
    [Fact]
    public void SplitMessage_WithValidSingleWordCallsigns_ReturnsCorrectTokens()
    {
        // Arrange
        var rawMessage = "ABC123 DEF456 Hello World";

        var communicators = new ActiveCommunicatorList();
        communicators.AddCommunicator(new Aircraft("ABC123", "Test"));
        communicators.AddCommunicator(new Aircraft("DEF456", "Test"));

        // Act
        var result = Message.SplitMessage(rawMessage, communicators);

        // Assert
        Assert.Equal(4, result.Count);

        Assert.True(result[0].IsCallsign);
        Assert.Equal("ABC123", result[0].Content);
        Assert.True(result[1].IsCallsign);
        Assert.Equal("DEF456", result[1].Content);
        Assert.False(result[2].IsCallsign);
        Assert.Equal("Hello", result[2].Content);
        Assert.False(result[3].IsCallsign);
        Assert.Equal("World", result[3].Content);
    }

    [Fact]
    public void SplitMessage_WithNoCallsigns_ReturnsAllContentTokens()
    {
        // Arrange
        var rawMessage = "Hello World Test";
        var communicators = new ActiveCommunicatorList();

        // Act
        var result = Message.SplitMessage(rawMessage, communicators);

        // Assert
        Assert.Equal(3, result.Count);
        Assert.All(result, token => Assert.False(token.IsCallsign));
        Assert.Equal("Hello", result[0].Content);
        Assert.Equal("World", result[1].Content);
        Assert.Equal("Test", result[2].Content);
    }

    [Fact]
    public void SplitMessage_WithEmptyMessage_ReturnsEmptyList()
    {
        // Arrange
        var rawMessage = "";
        var communicators = new ActiveCommunicatorList();

        // Act
        var result = Message.SplitMessage(rawMessage, communicators);

        // Assert
        Assert.Empty(result);
    }

    [Fact]
    public void SplitMessage_WithMultiWordCallsigns_ReturnsCorrectTokens()
    {
        // Arrange
        var rawMessage = "New York Center Delta Air Lines 123 request flight level 350";
        var communicators = new ActiveCommunicatorList();
        communicators.AddCommunicator(new Aircraft("New York Center", "ATC"));
        communicators.AddCommunicator(new Aircraft("Delta Air Lines 123", "Aircraft"));

        // Act
        var result = Message.SplitMessage(rawMessage, communicators);

        // Assert
        Assert.Equal(6, result.Count);

        Assert.True(result[0].IsCallsign);
        Assert.Equal("New York Center", result[0].Content);
        Assert.True(result[1].IsCallsign);
        Assert.Equal("Delta Air Lines 123", result[1].Content);
        Assert.False(result[2].IsCallsign);
        Assert.Equal("request", result[2].Content);
        Assert.False(result[3].IsCallsign);
        Assert.Equal("flight", result[3].Content);
        Assert.False(result[4].IsCallsign);
        Assert.Equal("level", result[4].Content);
        Assert.False(result[5].IsCallsign);
        Assert.Equal("350", result[5].Content);
    }
}
