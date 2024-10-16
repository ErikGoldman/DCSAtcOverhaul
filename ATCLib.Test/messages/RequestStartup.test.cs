using ATCLib.Messages;

namespace ATCLib.Test;

public class StartupMessageTests
{
    [Fact]
    public void Parse_WithRequestStartupMessage_ReturnsRequestStartupPayload()
    {
        // Arrange
        var rawMessage = "ABC123 requesting startup";
        var state = new ATCState();
        var abc123 = new Aircraft("ABC123", "Test");
        state.ActiveCommunicators.AddCommunicator(abc123);
        state.MessagePayloadTypes.Add(new RequestStartupMessagePayloadParser());

        // Act
        var result = Message.Parse(rawMessage, state);

        // Assert
        Assert.NotNull(result);
        Assert.Equal(rawMessage, result.RawMessage);
        Assert.Equal(abc123, result.From);
        Assert.Null(result.To);
        Assert.IsType<RequestStartup>(result.Payload);
        Assert.Empty(((RequestStartup)result.Payload).OtherPlanes);
    }

    [Fact]
    public void Parse_WithRequestStartupAndWeatherMessage_ReturnsRequestStartupPayload()
    {
        // Arrange
        var rawMessage = "DEF456 request startup and weather information";
        var state = new ATCState();
        var def456 = new Aircraft("DEF456", "Test");
        state.ActiveCommunicators.AddCommunicator(def456);
        state.MessagePayloadTypes.Add(new RequestStartupMessagePayloadParser());

        // Act
        var result = Message.Parse(rawMessage, state);

        // Assert
        Assert.NotNull(result);
        Assert.Equal(rawMessage, result.RawMessage);
        Assert.Equal(def456, result.From);
        Assert.Null(result.To);
        Assert.IsType<RequestStartup>(result.Payload);
        Assert.Empty(((RequestStartup)result.Payload).OtherPlanes);
    }

    [Fact]
    public void Parse_WithRequestStartupMessageAndOtherPlanes_ReturnsRequestStartupPayloadWithOtherPlanes()
    {
        // Arrange
        var rawMessage = "Tower GHI789 request startup with JKL012 and MNO345";
        var state = new ATCState();
        var tower = new Aircraft("Tower", "Test");
        var ghi789 = new Aircraft("GHI789", "Test");
        var jkl012 = new Aircraft("JKL012", "Test");
        var mno345 = new Aircraft("MNO345", "Test");
        state.ActiveCommunicators.AddCommunicator(tower);
        state.ActiveCommunicators.AddCommunicator(ghi789);
        state.ActiveCommunicators.AddCommunicator(jkl012);
        state.ActiveCommunicators.AddCommunicator(mno345);
        state.MessagePayloadTypes.Add(new RequestStartupMessagePayloadParser());

        // Act
        var result = Message.Parse(rawMessage, state);

        // Assert
        Assert.NotNull(result);
        Assert.Equal(rawMessage, result.RawMessage);
        Assert.Equal(ghi789, result.From);
        Assert.Equal(tower, result.To);
        Assert.IsType<RequestStartup>(result.Payload);
        var otherPlanes = ((RequestStartup)result.Payload).OtherPlanes;
        Assert.Equal(2, otherPlanes.Count);
        Assert.Contains(otherPlanes, p => p.Content == "JKL012");
        Assert.Contains(otherPlanes, p => p.Content == "MNO345");
    }

    [Fact]
    public void Parse_WithRequestStartupMessageAndOtherPlanes_DifferentPermutations()
    {
        var state = new ATCState();
        var tower = new Aircraft("Tower", "Test");
        var ghi789 = new Aircraft("GHI789", "Test");
        var jkl012 = new Aircraft("JKL012", "Test");
        var mno345 = new Aircraft("MNO345", "Test");

        // Arrange
        state.ActiveCommunicators.AddCommunicator(tower);
        state.ActiveCommunicators.AddCommunicator(ghi789);
        state.ActiveCommunicators.AddCommunicator(jkl012);
        state.ActiveCommunicators.AddCommunicator(mno345);
        state.MessagePayloadTypes.Add(new RequestStartupMessagePayloadParser());

        // Act
        var result = Message.Parse("Tower GHI789 requesting engine start", state);
        Assert.NotNull(result);
        Assert.IsType<RequestStartup>(result.Payload);

        result = Message.Parse("Tower GHI789 request engine powerup", state);
        Assert.NotNull(result);
        Assert.IsType<RequestStartup>(result.Payload);

        result = Message.Parse("Tower GHI789 request start with weather", state);
        Assert.NotNull(result);
        Assert.IsType<RequestStartup>(result.Payload);
    }

    [Fact]
    public void Parse_WithNonMatchingMessage_ReturnsNull()
    {
        // Arrange
        var rawMessage = "PQR678 hello";
        var state = new ATCState();
        var pqr678 = new Aircraft("PQR678", "Test");
        state.ActiveCommunicators.AddCommunicator(pqr678);
        state.MessagePayloadTypes.Add(new RequestStartupMessagePayloadParser());

        // Act
        var result = Message.Parse(rawMessage, state);

        // Assert
        Assert.NotNull(result);
        Assert.Equal(pqr678, result.From);
        Assert.IsType<EmptyMessagePayload>(result.Payload);
    }

    [Fact]
    public void Parse_WithRequestStartupMessageAndOtherPlanes_ReturnsRequestStartupPayloadWithOtherPlanes2()
    {
        // Arrange
        var rawMessage = "Tower GHI789 roger that flight of three with JKL012 and MNO345 requesting startup and weather info";
        var state = new ATCState();
        var tower = new Aircraft("Tower", "Test");
        var ghi789 = new Aircraft("GHI789", "Test");
        var jkl012 = new Aircraft("JKL012", "Test");
        var mno345 = new Aircraft("MNO345", "Test");
        state.ActiveCommunicators.AddCommunicator(tower);
        state.ActiveCommunicators.AddCommunicator(ghi789);
        state.ActiveCommunicators.AddCommunicator(jkl012);
        state.ActiveCommunicators.AddCommunicator(mno345);
        state.MessagePayloadTypes.Add(new RequestStartupMessagePayloadParser());

        // Act
        var result = Message.Parse(rawMessage, state);

        // Assert
        Assert.NotNull(result);
        Assert.Equal(rawMessage, result.RawMessage);
        Assert.Equal(ghi789, result.From);
        Assert.Equal(tower, result.To);
        Assert.IsType<RequestStartup>(result.Payload);
        var otherPlanes = ((RequestStartup)result.Payload).OtherPlanes;
        Assert.Equal(2, otherPlanes.Count);
        Assert.Contains(otherPlanes, p => p.Content == "JKL012");
        Assert.Contains(otherPlanes, p => p.Content == "MNO345");
    }
}
