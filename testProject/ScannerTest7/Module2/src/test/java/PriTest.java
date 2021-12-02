import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;

public class PriTest {
    Pri pri = new Pri();

    @Test
    public void testPriStm() {
        assertEquals(1, pri.prStm());
    }
}
