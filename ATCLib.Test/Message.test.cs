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

    [Fact]
    public void SplitMessage_WithAmbiguousPrefix_ReturnsCorrectTokens()
    {
        // Arrange
        var rawMessage = "Delta 123 Delta 456 request flight level 350";
        var communicators = new ActiveCommunicatorList();
        communicators.AddCommunicator(new Aircraft("Delta 123", "Aircraft"));
        communicators.AddCommunicator(new Aircraft("Delta 456", "Aircraft"));

        // Act
        var result = Message.SplitMessage(rawMessage, communicators);

        // Assert
        Assert.Equal(6, result.Count);

        Assert.True(result[0].IsCallsign);
        Assert.Equal("Delta 123", result[0].Content);
        Assert.True(result[1].IsCallsign);
        Assert.Equal("Delta 456", result[1].Content);
        Assert.False(result[2].IsCallsign);
        Assert.Equal("request", result[2].Content);
        Assert.False(result[3].IsCallsign);
        Assert.Equal("flight", result[3].Content);
        Assert.False(result[4].IsCallsign);
        Assert.Equal("level", result[4].Content);
        Assert.False(result[5].IsCallsign);
        Assert.Equal("350", result[5].Content);
    }

    [Fact]
    public void SplitMessage_WithAmbiguousPrefixNotCallsign_ReturnsCorrectTokens()
    {
        // Arrange
        var rawMessage = "Hey Delta Charlie how are you doing Delta 456 sup";
        var communicators = new ActiveCommunicatorList();
        communicators.AddCommunicator(new Aircraft("Delta Charlie 123", "Aircraft"));
        communicators.AddCommunicator(new Aircraft("Delta 456", "Aircraft"));

        // Act
        var result = Message.SplitMessage(rawMessage, communicators);

        // Assert
        Assert.Equal(9, result.Count);

        Assert.False(result[0].IsCallsign);
        Assert.Equal("Hey", result[0].Content);
        Assert.False(result[1].IsCallsign);
        Assert.Equal("Delta", result[1].Content);
        Assert.False(result[2].IsCallsign);
        Assert.Equal("Charlie", result[2].Content);
        Assert.False(result[3].IsCallsign);
        Assert.Equal("how", result[3].Content);
        Assert.False(result[4].IsCallsign);
        Assert.Equal("are", result[4].Content);
        Assert.False(result[5].IsCallsign);
        Assert.Equal("you", result[5].Content);
        Assert.False(result[6].IsCallsign);
        Assert.Equal("doing", result[6].Content);
        Assert.True(result[7].IsCallsign);
        Assert.Equal("Delta 456", result[7].Content);
        Assert.False(result[8].IsCallsign);
        Assert.Equal("sup", result[8].Content);
    }

}
