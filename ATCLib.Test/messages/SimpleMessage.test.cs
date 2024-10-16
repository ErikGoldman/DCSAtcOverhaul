namespace ATCLib.Test;

using ATCLib.Messages;
using Xunit;

public class SimpleMessageTests
{
    [Fact]
    public void Parse_WithRogerMessage_ReturnsRogerPayload()
    {
        // Arrange
        var rawMessage = "ABC123 roger";
        var state = new ATCState();
        var abc123 = new Aircraft("ABC123", "Test");
        state.ActiveCommunicators.AddCommunicator(abc123);
        state.MessagePayloadTypes.Add(new RogerMessagePayloadParser());

        // Act
        var result = Message.Parse(rawMessage, state);

        // Assert
        Assert.NotNull(result);
        Assert.Equal(rawMessage, result.RawMessage);
        Assert.Equal(abc123, result.From);
        Assert.Null(result.To);
        Assert.IsType<Roger>(result.Payload);
        Assert.Equal("roger", result.Payload.RawContent);
    }

    [Fact]
    public void Parse_WithAffirmativeMessage_ReturnsRogerPayload()
    {
        // Arrange
        var rawMessage = "DEF456 affirmative";
        var state = new ATCState();
        var def456 = new Aircraft("DEF456", "Test");
        state.ActiveCommunicators.AddCommunicator(def456);
        state.MessagePayloadTypes.Add(new RogerMessagePayloadParser());

        // Act
        var result = Message.Parse(rawMessage, state);

        // Assert
        Assert.NotNull(result);
        Assert.Equal(rawMessage, result.RawMessage);
        Assert.Equal(def456, result.From);
        Assert.Null(result.To);
        Assert.IsType<Roger>(result.Payload);
        Assert.Equal("affirmative", result.Payload.RawContent);
    }

    [Fact]
    public void Parse_WithCopyMessage_ReturnsRogerPayload()
    {
        // Arrange
        var rawMessage = "GHI789 copy";
        var state = new ATCState();
        var ghi789 = new Aircraft("GHI789", "Test");
        state.ActiveCommunicators.AddCommunicator(ghi789);
        state.MessagePayloadTypes.Add(new RogerMessagePayloadParser());

        // Act
        var result = Message.Parse(rawMessage, state);

        // Assert
        Assert.NotNull(result);
        Assert.Equal(rawMessage, result.RawMessage);
        Assert.Equal(ghi789, result.From);
        Assert.Null(result.To);
        Assert.IsType<Roger>(result.Payload);
        Assert.Equal("copy", result.Payload.RawContent);
    }

    [Fact]
    public void Parse_WithCopyMessage_ReturnsRogerPayload2()
    {
        // Arrange
        var rawMessage = "Copy that GHI789";
        var state = new ATCState();
        var ghi789 = new Aircraft("GHI789", "Test");
        state.ActiveCommunicators.AddCommunicator(ghi789);
        state.MessagePayloadTypes.Add(new RogerMessagePayloadParser());

        // Act
        var result = Message.Parse(rawMessage, state);

        // Assert
        Assert.NotNull(result);
        Assert.Equal(rawMessage, result.RawMessage);
        Assert.Equal(ghi789, result.From);
        Assert.Null(result.To);
        Assert.IsType<Roger>(result.Payload);
    }

    [Fact]
    public void Parse_WithCopyMessage_ReturnsEmpty()
    {
        // Arrange
        var rawMessage = "Ok I got it GHI789";
        var state = new ATCState();
        var ghi789 = new Aircraft("GHI789", "Test");
        state.ActiveCommunicators.AddCommunicator(ghi789);
        state.MessagePayloadTypes.Add(new RogerMessagePayloadParser());

        // Act
        var result = Message.Parse(rawMessage, state);

        // Assert
        Assert.NotNull(result);
        Assert.Equal(rawMessage, result.RawMessage);
        Assert.Equal(ghi789, result.From);
        Assert.Null(result.To);
        Assert.IsType<EmptyMessagePayload>(result.Payload);
    }

    [Fact]
    public void Parse_WithCopyMessage_ReturnsRogerPayload3()
    {
        // Arrange
        var rawMessage = "GHI789";
        var state = new ATCState();
        var ghi789 = new Aircraft("GHI789", "Test");
        state.ActiveCommunicators.AddCommunicator(ghi789);
        state.MessagePayloadTypes.Add(new RogerMessagePayloadParser());

        // Act
        var result = Message.Parse(rawMessage, state);

        // Assert
        Assert.NotNull(result);
        Assert.Equal(rawMessage, result.RawMessage);
        Assert.Equal(ghi789, result.From);
        Assert.Null(result.To);
        Assert.IsType<Roger>(result.Payload);
    }

    [Fact]
    public void Parse_WithSayAgainMessage_ReturnsSayAgainPayload()
    {
        // Arrange
        var rawMessage = "JKL012 say again";
        var state = new ATCState();
        var jkl012 = new Aircraft("JKL012", "Test");
        state.ActiveCommunicators.AddCommunicator(jkl012);
        state.MessagePayloadTypes.Add(new SayAgainMessagePayloadParser());

        // Act
        var result = Message.Parse(rawMessage, state);

        // Assert
        Assert.NotNull(result);
        Assert.Equal(rawMessage, result.RawMessage);
        Assert.Equal(jkl012, result.From);
        Assert.Null(result.To);
        Assert.IsType<SayAgain>(result.Payload);
        Assert.Equal("say again", result.Payload.RawContent);
    }

    [Fact]
    public void Parse_WithRepeatMessage_ReturnsSayAgainPayload()
    {
        // Arrange
        var rawMessage = "MNO345 repeat";
        var state = new ATCState();
        var mno345 = new Aircraft("MNO345", "Test");
        state.ActiveCommunicators.AddCommunicator(mno345);
        state.MessagePayloadTypes.Add(new SayAgainMessagePayloadParser());

        // Act
        var result = Message.Parse(rawMessage, state);

        // Assert
        Assert.NotNull(result);
        Assert.Equal(rawMessage, result.RawMessage);
        Assert.Equal(mno345, result.From);
        Assert.Null(result.To);
        Assert.IsType<SayAgain>(result.Payload);
        Assert.Equal("repeat", result.Payload.RawContent);
    }

    [Fact]
    public void Parse_WithRepeatMessage_ReturnsSayAgainPayload2()
    {
        // Arrange
        var rawMessage = "control tower say that again please MNO345";
        var state = new ATCState();
        var mno345 = new Aircraft("MNO345", "Test");
        var other = new Aircraft("control tower", "Test");
        state.ActiveCommunicators.AddCommunicator(mno345);
        state.ActiveCommunicators.AddCommunicator(other);
        state.MessagePayloadTypes.Add(new SayAgainMessagePayloadParser());

        // Act
        var result = Message.Parse(rawMessage, state);

        // Assert
        Assert.NotNull(result);
        Assert.Equal(rawMessage, result.RawMessage);
        Assert.Equal(mno345, result.From);
        Assert.Equal(other, result.To);
        Assert.IsType<SayAgain>(result.Payload);
    }

    [Fact]
    public void Parse_WithNonMatchingMessage_ReturnsNull()
    {
        // Arrange
        var rawMessage = "PQR678 hello";
        var state = new ATCState();
        var pqr678 = new Aircraft("PQR678", "Test");
        state.ActiveCommunicators.AddCommunicator(pqr678);
        state.MessagePayloadTypes.Add(new RogerMessagePayloadParser());
        state.MessagePayloadTypes.Add(new SayAgainMessagePayloadParser());

        // Act
        var result = Message.Parse(rawMessage, state);

        // Assert
        Assert.NotNull(result);
        Assert.Equal(pqr678, result.From);
        Assert.IsType<EmptyMessagePayload>(result.Payload);
    }
}